package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

const (
	index       = "index.html"
	layout      = "layout.html"
	index_path  = "templates/index.html"
	layout_path = "templates/layout.html"
	footer_path = "templates/footer.html"
	head_path   = "templates/head.html"
	url         = "DB_URL"
	home_page   = "HOME_PAGE"
)

var tmpl = template.Must(template.ParseFiles(index_path, layout_path, head_path, footer_path))
var todaysRecipe int
var maxOffset int

func main() {
	fmt.Println("Running main")
	getMaxRows()
	todaysRecipe = rand.Intn(maxOffset)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	c := cron.New()
	c.AddFunc("00 00 * * *", func() {
		todaysRecipe = rand.Intn(maxOffset)
	})

	c.Start()

	http.HandleFunc("/", handler)
	// http.HandleFunc("/execute", handlerExecute)
	fmt.Println("Starting service on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	recipe := getRecipe()
	err := tmpl.ExecuteTemplate(w, index,
		struct {
			Recipe    Recipe
			Home_page string
		}{recipe, os.Getenv(home_page)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func getRecipe() Recipe {
	conn, err := pgx.Connect(context.Background(), os.Getenv(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var recipe Recipe
	fmt.Println(todaysRecipe)
	fmt.Println("select name, minutes, submitted, tags, nutrition, n_steps,steps, description, ingredients, n_ingredients from recipes order by id limit 1 offset " + strconv.Itoa(todaysRecipe) + ";")
	err = conn.QueryRow(context.Background(), "select name, minutes, submitted, tags, nutrition, n_steps,steps, description, ingredients, n_ingredients from recipes order by id limit 1 offset "+strconv.Itoa(todaysRecipe)+";").
		Scan(
			&recipe.Name, &recipe.Minutes, &recipe.Submitted, &recipe.Tags,
			&recipe.Nutrition, &recipe.N_steps, &recipe.Steps, &recipe.Description, &recipe.Ingredients, &recipe.N_ingredients,
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	return recipe
}

func getMaxRows() {
	conn, err := pgx.Connect(context.Background(), os.Getenv(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select count(*) from recipes;").
		Scan(&maxOffset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get my query off: %v\n", err)
		os.Exit(1)
	}
}
