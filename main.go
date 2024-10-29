package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

const (
	index          = "index.html"
	pictures       = "pictures.html"
	index_path     = "templates/index.html"
	footer_path    = "templates/footer.html"
	head_path      = "templates/head.html"
	header_path    = "templates/header.html"
	pictures_path  = "templates/pictures.html"
	url            = "DB_URL"
	pictures_url   = "/pictures"
	home_page      = "HOME_PAGE"
	pictures_page  = "PICTURES_PAGE"
	picture_upload = "PICTURES_UPLOAD"
)

var tmpl = template.Must(template.ParseFiles(
	index_path, head_path, footer_path, pictures_path, header_path),
)
var todaysRecipe int
var maxOffset int

func main() {
	fmt.Println("Running main")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	getMaxRows(ctx)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fs = http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	c := cron.New()
	c.AddFunc("00 00 * * *", func() {
		todaysRecipe = rand.Intn(maxOffset)
	})

	c.Start()

	http.HandleFunc("/", handler)
	http.HandleFunc("/pictures", handlerPictures)
	http.HandleFunc("/save_picture", handlerSavePicture)

	fmt.Println("Starting service on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	recipe := getRecipe(r.Context())
	err := tmpl.ExecuteTemplate(w, index,
		struct {
			Recipe        Recipe
			Home_page     string
			Pictures_page string
			Path          string
		}{
			recipe,
			os.Getenv(home_page),
			os.Getenv(pictures_page),
			r.URL.Path,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerPictures(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb(r.Context())
	defer conn.Close(r.Context())

	fmt.Println("PATH: ", r.URL.Path)
	var paths []string
	rows, err := conn.Query(r.Context(), "select path from pictures;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	paths, err = pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		http.Error(w, "Collecting rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, pictures, struct {
		Paths          []string
		Home_page      string
		Pictures_page  string
		Picture_upload string
		Path           string
	}{
		paths,
		os.Getenv(home_page),
		os.Getenv(pictures_page),
		picture_upload,
		r.URL.Path,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerSavePicture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This is a POST method you pirate!", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	conn := connectToDb(r.Context())
	defer conn.Close(r.Context())

	cmd, err := conn.Exec(
		r.Context(),
		"INSERT INTO pictures (user_id, path) VALUES ($1, $2)", 1, ".",
	)
	if err != nil {
		http.Error(w, "Insert err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cmd.Insert()

	var fileId int
	err = conn.QueryRow(r.Context(), "SELECT id FROM pictures WHERE user_id="+"1"+" ORDER BY id DESC LIMIT 1").Scan(&fileId)

	fileName := strconv.Itoa(fileId)

	formFile, _, err := r.FormFile(picture_upload)
	if err != nil {
		http.Error(w, "Form file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer formFile.Close()

	filePath := filepath.Join("./public/pictures", fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Creation err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, formFile); err != nil {
		http.Error(w, "Copy error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cmd, err = conn.Exec(r.Context(), "UPDATE pictures SET (path, updated_at) = ($1, CURRENT_TIMESTAMP) where id=$2", filePath, fileName)
	if err != nil {
		http.Error(w, "Insert err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cmd.Update()

	http.Redirect(w, r, pictures_url, http.StatusSeeOther)
}

type Recipe struct {
	Name          string
	Minutes       int
	Submitted     time.Time
	Tags          []string
	Nutrition     string
	N_steps       int
	Steps         []string
	Description   string
	Ingredients   []string
	N_ingredients int
}

func getRecipe(ctx context.Context) Recipe {
	conn := connectToDb(ctx)
	defer conn.Close(ctx)

	var recipe Recipe
	fmt.Println(todaysRecipe)
	fmt.Println("select name, minutes, submitted, tags, nutrition, n_steps,steps, description, ingredients, n_ingredients from recipes order by id limit 1 offset " + strconv.Itoa(todaysRecipe) + ";")
	err := conn.QueryRow(ctx, "select name, minutes, submitted, tags, nutrition, n_steps,steps, description, ingredients, n_ingredients from recipes order by id limit 1 offset "+strconv.Itoa(todaysRecipe)+";").
		Scan(
			&recipe.Name, &recipe.Minutes, &recipe.Submitted, &recipe.Tags,
			&recipe.Nutrition, &recipe.N_steps, &recipe.Steps, &recipe.Description, &recipe.Ingredients, &recipe.N_ingredients,
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	return recipe
}

func getMaxRows(ctx context.Context) {
	conn := connectToDb(ctx)
	defer conn.Close(ctx)

	err := conn.QueryRow(ctx, "select recipes_offset from dailies where date(created_at) = CURRENT_DATE;").
		Scan(&todaysRecipe)
	fmt.Println(todaysRecipe)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Dailies: Unable to get my query off: %v\n", err)
	}

	err = conn.QueryRow(ctx, "select count(id) from recipes;").
		Scan(&maxOffset)

	if todaysRecipe == 0 {
		todaysRecipe = rand.Intn(maxOffset)
		cmd, err := conn.Exec(ctx, "INSERT INTO dailies (recipes_offset) VALUES ($1)", todaysRecipe)
		cmd.Insert()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't insert: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = conn.QueryRow(ctx, "select recipes_offset from dailies;").
			Scan(&todaysRecipe)
	}

}

func connectToDb(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
