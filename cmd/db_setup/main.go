package main

import (
	"database/sql"
	"log"
	"os"

	las "github.com/dmholtz/logic_app_server"

	_ "github.com/mattn/go-sqlite3"
)

// This tool creates the sqlite database file from the schema.sql file.
//
// It useses golang's SQL package and the sqlite3 driver.
// Source: https://pkg.go.dev/github.com/mattn/go-sqlite3?utm_source=godoc
func main() {
	CreateDatabase()
	InsertAchievements()
	AddDummyData()
}

func CreateDatabase() {
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

func InsertAchievements() {
	// read the SQL definition of achievements
	bytes, err := os.ReadFile("db/achievements.sql")
	if err != nil {
		log.Fatal(err)
		return
	}
	sqlAchievementScript := string(bytes)

	// open the database file
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// execute the SQL script to insert the achievements
	_, err = db.Exec(sqlAchievementScript)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Default().Println("Achievements inserted successfully")
}

func AddDummyData() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	userStore := &las.MyUserStore{DB: db}

	// insert dummy players
	userStore.Signup(las.Credentials{Username: "user1", Password: "user1"})
	userStore.Signup(las.Credentials{Username: "user2", Password: "user2"})
	userStore.Signup(las.Credentials{Username: "user3", Password: "user3"})

	// login dummy players
	userStore.Login(las.Credentials{Username: "user1", Password: "user1"})

	// insert dummy quiz participations
	db.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, points) VALUES (1,1,1,10)")
	db.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, points) VALUES (2,1,1,5)")
	db.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, points) VALUES (3,1,1,10)")
	db.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, points) VALUES (1,2,1,2)")
	db.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, points) VALUES (2,2,1,2)")

	// insert dummy achieved relations
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (1,1)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (1,2)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (1,4)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (1,5)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (1,6)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (2,1)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (3,1)")
	db.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (3,2)")
}
