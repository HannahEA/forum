package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

type userDetails struct {
	ID       int
	email    string
	username string
	password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func newUser(email, username, password string, db *sql.DB) {
	hash, err := HashPassword(password)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, errNewUser := db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hash)
	if errNewUser != nil {
		fmt.Printf("The error is %v", errNewUser.Error())
		log.Fatal()
	}
}

//userExsists checks if the username entered is already taken. If it is the function returns true.
func userExist(email, username string, db *sql.DB) (bool, string) {
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
		log.Fatal(err1.Error())
	}

	count1 := 0
	for rows1.Next() {
		count1++
	}

	if count == 0 && count1 == 0 {
		return true, "You have successfully registered"
	} else if count1 == 1 && count == 1 {
		return false, "This account already exists. Please sign in"

	} else if count == 1 {
		return false, "This email is already taken. Please try a different email address"
	} else {
		return false, "This username is already taken. Please try a different username"
	}
}

func ValidEmail(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatal(err.Error())
	}
	count := 0

	for rows.Next() {
		count++
	}

	if count == 0 {
		return false
	} else {
		return true
	}
}

func LoginValidator(email, password string, db *sql.DB) bool {
	rows1, err1 := db.Query("SELECT ID, email, username, password FROM users WHERE email = ?", email)

	if err1 != nil {
		log.Fatal(err1.Error())
	}

	var u userDetails

	err := rows1.Scan(
		&u.ID,
		&u.email,
		&u.username,
		&u.password,
	)

	if err != nil {
		fmt.Println("SCANNING ERROR")
		log.Fatal(err.Error())
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))

	return hashErr == nil

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/login.html"))
	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/register.html"))
	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}

func registration2(w http.ResponseWriter, r *http.Request) {

	userN := r.FormValue("username")
	email := r.FormValue("email")
	pass := r.FormValue("password")

	exist, value := userExist(email, userN, sqliteDatabase)

	tpl := template.Must(template.ParseGlob("templates/register2.html"))

	if !exist {
		if err := tpl.Execute(w, value); err != nil {
			log.Fatal(err.Error())
		}
	} else {

		newUser(email, userN, pass, sqliteDatabase)

		if err := tpl.Execute(w, value); err != nil {
			log.Fatal(err.Error())
		}

	}

}
