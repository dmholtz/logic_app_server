package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// This tool creates the sqlite database file from the schema.sql file.
// It useses golang's SQL package and the sqlite3 driver.
// Source: https://pkg.go.dev/github.com/mattn/go-sqlite3?utm_source=godoc
func main() {
	// read the SQL init script to create the schema
	bytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlInitScript := string(bytes)

	// remove any existing database file
	_ = os.Remove("./db.sqlite3")

	// create the database file
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// execute the SQL init script to create the tables according to the schema
	_, err = db.Exec(sqlInitScript)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Default().Println("Database created successfully")
}
