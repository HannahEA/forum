package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB
var Person userDetails

func main() {

	//Open the database SQLite file and create the database table
	database, err1 := sql.Open("sqlite3", "sqlite-database.db")
	sqliteDatabase = database

	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	defer sqliteDatabase.Close()

	http.HandleFunc("/", Home)
	http.HandleFunc("/log", LoginHandler)
	http.HandleFunc("/login", LoginResult)
	http.HandleFunc("/register", registration)
	http.HandleFunc("/registration", registration2)
	http.HandleFunc("/logout", LogOut)
	http.ListenAndServe(":8080", nil)

}
