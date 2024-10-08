package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type query struct {
	name           string
	id             int
	minutes        int
	contributor_id int
	submitted      string
	tags           string
	nutrition      string
	n_steps        int
	steps          string
	description    string
	ingredients    string
	n_ingredients  int
}

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://dima:dima@localhost:5433/food")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	file, err := os.Open("./small.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var query query

	for scanner.Scan() {
		line := scanner.Text()
		query.name, line = parse(line)
		query.id, line = parseNumber(line)
		query.minutes, line = parseNumber(line)
		query.contributor_id, line = parseNumber(line)
		query.submitted, line = parse(line)
		query.tags, line = parseArray(line)
		query.nutrition, line = parseArray(line)
		query.n_steps, line = parseNumber(line)
		query.steps, line = parseArray(line)
		query.description, line = parse(line)
		query.ingredients, line = parseArray(line)
		query.n_ingredients, _ = strconv.Atoi(line)
	}

	fmt.Println(query)
}

func parseNumber(line string) (int, string) {
	comma := strings.Index(line, ",")
	num, err := strconv.Atoi(line[:comma])
	if err != nil {
		log.Fatal("Not a numba", err)
	}
	return num, line[comma+1:]
}

func parse(line string) (string, string) {
	comma := strings.Index(line, ",")
	return line[:comma], line[comma+1:]
}

func parseArray(line string) (string, string) {
	comma := strings.Index(line, "]") + 2
	array := line[:comma]
	array = strings.Replace(array, "[", "{", 1)
	array = strings.Replace(array, "]", "}", 1)
	array = strings.ReplaceAll(array, "'", "\"")
	// array = "'" + array[1:]
	// array = "'" + array[:len(array)-2]
	return array, line[comma+1:]
}
