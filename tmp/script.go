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

const (
	EMPTY_LINE = -1
	TOO_SHORT  = 200
	SKIP_COMMA = 3
	EMPTY_DESC = 2
)

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
	var full_statement string

	for scanner.Scan() {
		line := scanner.Text()
		full_statement += line
		idx := strings.LastIndex(full_statement, ",")

		if idx == EMPTY_LINE {
			fmt.Println("There is no comma!")
			continue
		}

		if len(line) == TOO_SHORT {
			fmt.Println("too short")
			continue
		}

		n_steps, err := strconv.Atoi(full_statement[idx+1:])

		if err != nil {
			fmt.Printf("This is not a number, full statement: %v\n\n", full_statement[idx+1:])
			continue
		}

		line = full_statement
		full_statement = ""
		fmt.Printf("Congrats you have a full line, here's the steps: %v, len: %v\n\n", n_steps, len(line))

		fmt.Println("name")
		query.name, line = parse(line)
		fmt.Println("\n\nid")
		query.id, line = parseNumber(line)
		fmt.Println("\n\nminutes")
		query.minutes, line = parseNumber(line)
		fmt.Println("\n\ncontributor_id")
		query.contributor_id, line = parseNumber(line)
		fmt.Println("\n\nsubmitted")
		query.submitted, line = parse(line)
		fmt.Println("\n\ntags")
		query.tags, line = parseArray(line)
		fmt.Println("\n\nnutrition")
		query.nutrition, line = parseArray(line)
		fmt.Println("\n\nn_steps")
		query.n_steps, line = parseNumber(line)
		fmt.Println("\n\nsteps")
		query.steps, line = parseArray(line)
		fmt.Println("\n\ndescription")
		query.description, line = parseDescription(line)
		fmt.Println("\n\ningredients")
		query.ingredients, line = parseArray(line)
		fmt.Println("\n\nn_ingredients")
		query.n_ingredients, err = strconv.Atoi(line)

		_, err = conn.Exec(context.Background(), query.getString(),
			query.name, query.id, query.minutes, query.contributor_id,
			query.submitted, query.tags, query.nutrition, query.n_steps, query.steps,
			query.description, query.ingredients, query.n_ingredients,
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(query)
}

func parseNumber(line string) (int, string) {
	comma := strings.Index(line, ",")
	value := line[:comma]
	fmt.Printf("\nI'm parsing a NUMBER in this string value: %v", value)
	num, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("\nNot a numba", err)
	}
	fmt.Print("\nline: ", line)
	return num, line[comma+1:]
}

func parse(line string) (string, string) {
	comma := strings.Index(line, ",")
	value := line[:comma]
	fmt.Printf("I'm parsing a STRING in this string value: %v", value)
	return value, line[comma+1:]
}

func parseDescription(line string) (string, string) {
	comma := strings.Index(line, "[")
	if comma <= EMPTY_DESC {
		return "", line[comma-1:]
	}
	value := line[1 : comma-3]
	fmt.Printf("\n\nI'm parsing a DESC in this string value: %v\n\n", value)
	return value, line[comma-1:]
}

func parseArray(line string) (string, string) {
	comma := strings.Index(line, "]")
	array := line[:comma+2]
	fmt.Printf("I'm parsing a ARRAY in this string value: %v", array)

	array = strings.Replace(array, "[", "{", 1)
	array = strings.Replace(array, "]", "}", 1)
	array = array[1:]
	array = array[:len(array)-1]

	return array, line[comma+SKIP_COMMA:]
}
