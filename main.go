package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func main() {

	//Open the database SQLite file and create the database table
	database, err1 := sql.Open("sqlite3", "sqlite-database.db")
	sqliteDatabase = database

	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	defer sqliteDatabase.Close()

	fmt.Println(LoginValidator("helld@hotmail.com", "rew", sqliteDatabase))
	fmt.Println(LoginValidator("helld@hotil.com", "HelloW", sqliteDatabase))

	http.HandleFunc("/", LoginHandler)
	http.HandleFunc("/register", registration)
	http.HandleFunc("/registration", registration2)
	http.ListenAndServe(":8080", nil)

}
