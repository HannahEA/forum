package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func newUser(email, username, password string, db *sql.DB) {
	_, errNewUser := db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, password)
	if errNewUser != nil {
		fmt.Printf("The error is %v", errNewUser.Error())
		log.Fatal()
	}

	// _, preAdded := add.Exec(email, username, password)
	// if preAdded != nil {
	// 	log.Fatal(preAdded.Error())
	// }
}

//userExsists checks if the username entered is already taken. If it is the function returns true.
func userExist(email, username string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatal(err.Error())
	}
	count := 0

	for rows.Next() {
		count++
	}

	rows1, err1 := db.Query("SELECT username FROM users WHERE username = ?", username)
	if err1 != nil {
		log.Fatal(err.Error())
	}

	count1 := 0
	for rows1.Next() {
		count1++
	}

	if count1 == 0 && count == 0 {
		return false
	} else {
		return true
	}
}
func registration(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/register.html"))
	if err := tpl.Execute(w,"HEllo"); err != nil {
		log.Fatal(err.Error())
	}

}
