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

func (q query) getString() string {

	// 	insert := `
	// insert into
	//   recipes (
	//     name,
	//     id,
	//     minutes,
	//     contributor_id,
	//     submitted,
	//     tags,
	//     nutrition,
	//     n_steps,
	//     steps,
	//     description,
	//     ingredients,
	//     n_ingredients
	//   )
	// values
	//   (
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v,
	//     %v
	//   );
	// `

	insert := `
insert into
  recipes (
    name,
    id,
    minutes,
    contributor_id,
    submitted,
    tags,
    nutrition,
    n_steps,
    steps,
    description,
    ingredients,
    n_ingredients
  )
values
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12
  );
`

	// output := fmt.Sprintf(
	// 	insert, "'"+q.name+"'", q.id, q.minutes, q.contributor_id,
	// 	"'"+q.submitted+"'", q.tags, q.nutrition, q.n_steps, q.steps,
	// 	"'"+q.description+"'", q.ingredients, q.n_ingredients)
	// fmt.Println(output)

	return insert
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
	line_numba := 1

	for scanner.Scan() {
		line := scanner.Text()
		line_numba++
		if len(line) < 2 {
			continue
		}
		fmt.Println("line ", line_numba, " :", line[len(line)-2:])
		_, err := strconv.Atoi(line[len(line)-2:])
		if err != nil {
			continue
		}
		query.name, line, err = parse(line)
		if err != nil {
			fmt.Println(err)
		}
		query.id, line, err = parseNumber(line)
		if err != nil {
			fmt.Println(err)
		}
		query.minutes, line, err = parseNumber(line)
		if err != nil {
			fmt.Println(err)
		}
		query.contributor_id, line, err = parseNumber(line)
		if err != nil {
			fmt.Println(err)
		}
		query.submitted, line, err = parse(line)
		if err != nil {
			fmt.Println(err)
		}
		query.tags, line, err = parseArray(line)
		if err != nil {
			fmt.Println("array in shambles")
		}
		query.nutrition, line, err = parseArray(line)
		if err != nil {
			fmt.Println("array in shambles")
		}
		query.n_steps, line, err = parseNumber(line)
		if err != nil {
			fmt.Println(err)
		}
		query.steps, line, err = parseArray(line)
		if err != nil {
			fmt.Println("array in shambles")
		}
		query.description, line, err = parseDescription(line)
		if err != nil {
			continue
		}
		query.ingredients, line, err = parseArray(line)
		if err != nil {
			fmt.Println("array in shambles")
		}
		query.n_ingredients, err = strconv.Atoi(line)
		if err != nil {
			fmt.Println("no numnba")
		}
		_, err = conn.Exec(context.Background(), query.getString(),
			query.name, query.id, query.minutes, query.contributor_id,
			query.submitted, query.tags, query.nutrition, query.n_steps, query.steps,
			query.description, query.ingredients, query.n_ingredients,
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	// fmt.Println(query)
}

func parseNumber(line string) (int, string, error) {
	comma := strings.Index(line, ",")
	if comma == -1 {
		return 0, "", fmt.Errorf("nope")
	}
	num, err := strconv.Atoi(line[:comma])
	if err != nil {
		fmt.Println("Not a numba", err)
		return 0, "", fmt.Errorf("numbers imma right?")
	}
	return num, line[comma+1:], nil
}

func parse(line string) (string, string, error) {
	comma := strings.Index(line, ",")
	if comma == -1 {
		return "", "", fmt.Errorf("guess nothing")
	}
	return line[:comma], line[comma+1:], nil
}

func parseDescription(line string) (string, string, error) {
	comma := strings.Index(line, "[")
	if comma < 1 {
		return "", "", fmt.Errorf("not 0")
	}
	return line[1 : comma-3], line[comma-1:], nil
}

func parseArray(line string) (string, string, error) {
	comma := strings.Index(line, "]")
	if comma == -1 {
		return "", "", fmt.Errorf("nothing")
	}
	array := line[:comma+2]
	array = strings.Replace(array, "[", "{", 1)
	array = strings.Replace(array, "]", "}", 1)
	// array = strings.ReplaceAll(array, "'", "\"")
	// array = "'" + array[1:]
	// array = array[:len(array)-1] + "'"
	// fmt.Println("array: ", array)
	if len(array)-1 <= 0 {
		return "", "", fmt.Errorf("asdfasdf")
	}
	array = array[1:]
	array = array[:len(array)-1]
	return array, line[comma+1:], nil
}
