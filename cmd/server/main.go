package main

import (
	"database/sql"
	"log"
	"net/http"

	las "github.com/dmholtz/logic_app_server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, dbErr := sql.Open("sqlite3", "./db.sqlite3")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	server := las.NewLogicAppServer(db)
	err := http.ListenAndServeTLS(":443", "cert/dev.cert.pem", "cert/dev.key.pem", server)
	//err := http.ListenAndServe(":443", server)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
