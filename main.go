package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//Open the database SQLite file and create the database table
	sqliteDatabase, err1 := sql.Open("sqlite3", "database/sqlite-database.db")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	defer sqliteDatabase.Close()

	http.HandleFunc("/", registration)
	http.ListenAndServe(":8080", nil)

	newUser("ndr@hotmail.com", "nater68", "HelloW", sqliteDatabase)
}
