package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

const (
	index = "index.html"
	url   = "DB_URL"
)

var tmpl = template.Must(template.ParseFiles(index))

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
	Name string
}

func getRecipe() Recipe {
	conn, err := pgx.Connect(context.Background(), os.Getenv(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var recipe Recipe
	err = conn.QueryRow(context.Background(), "select name from recipes limit 1;").Scan(&recipe.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return recipe
}
