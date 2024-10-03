package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user= dbname=food sslmode=disable password="
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	row, err := db.Query("select name from recipes limit 1")
	if err != nil {
		log.Fatal("Execute: ", err)
	}
	row.Scan()
	fmt.Println(row.Next())
	fmt.Println(row)
}
