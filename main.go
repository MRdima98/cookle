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
	index = "index.html"
	url   = "DB_URL"
)

var tmpl = template.Must(template.ParseFiles(index))
var todaysRecipe int
var maxOffset int

func main() {
	getMaxRows()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 1m", func() { todaysRecipe = rand.Intn(maxOffset) })
	c.Start()

	http.HandleFunc("/", handler)
	// http.HandleFunc("/execute", handlerExecute)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	recipe := getRecipe()
	err := tmpl.ExecuteTemplate(w, index, recipe)
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
	err = conn.QueryRow(context.Background(), "select name, minutes, submitted, tags, nutrition, n_steps,steps, description, ingredients, n_ingredients from recipes limit 1 offset "+strconv.Itoa(todaysRecipe)+";").
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
}
