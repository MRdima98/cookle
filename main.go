package main

import (
	"context"
	"fmt"
	"html/template"

	// "html/template"
	"log"
	"net/http"

	"os"

	"github.com/jackc/pgx/v5"
)

const (
	index = "index.html"
)

var tmpl = template.Must(template.ParseFiles(index))

func main() {
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
	urlExample := "postgres://dima:dima@localhost:5432/food"
	conn, err := pgx.Connect(context.Background(), urlExample)
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
