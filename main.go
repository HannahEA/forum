package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

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

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs)) // handling the CSS

	http.HandleFunc("/", Home)
	http.HandleFunc("/log", LoginHandler)
	http.HandleFunc("/login", LoginResult)
	http.HandleFunc("/register", registration)
	http.HandleFunc("/registration", registration2)
	http.HandleFunc("/logout", LogOut)
	http.HandleFunc("/new-post", Post)
	http.HandleFunc("/post-added", postAdded)
	http.ListenAndServe(":8080", nil)

}

func Executer(w http.ResponseWriter, file string) {
	tpl := template.Must(template.ParseGlob(file))

	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}
