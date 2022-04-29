package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//Create the SQLite database
	file, err := os.Create("sqlite-database.db")

	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	fmt.Println("SQL Databasefile created")

	//Open the database SQLite file and create the database table
	sqliteDatabase, err1 := sql.Open("sqlite3", "sqlite-database.db")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	defer sqliteDatabase.Close()

	//Create the database for each user
	usersTbl, errTbl := sqliteDatabase.Prepare(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"email" 	TEXT UNIQUE,
			"username"	TEXT,
			"password"	TEXT 
		);
	`)

	if errTbl != nil {
		fmt.Println("TABLE ERROR")
		log.Fatal(errTbl.Error())
	}

	usersTbl.Exec()

	//Create the posts table
	postsTbl, errPosts := sqliteDatabase.Prepare(`
	CREATE TABLE IF NOT EXSISTS "posts" (
		"postNum"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"userID"	INTEGER,
		"category"	TEXT,
		"likes" INTEGER,
		"dislikes" INTEGER,
		);
	`)

	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}

	postsTbl.Exec()

}

//newUser adds a new account with a unique username and password to the database.
func newUser(email, username, password string, db *sql.DB) {
	add, errNewUser := db.Prepare("INSERT INTO users (email, username, password) VALUES (?, ?, ?)")
	if errNewUser != nil {
		log.Fatal()
	}

	_, preAdded := add.Exec(email, username, password)
	if preAdded != nil {
		log.Fatal(preAdded.Error())
	}
}

//newPost creates a new post by a registered user
func newPost(userID int, category string, db *sql.DB) {
	add, err := db.Prepare("INSERT INTO posts (userID, category, likes, dislikes) VALUES (?,?, 0,0)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, errPost := add.Exec(userID, category)
	if errPost != nil {
		log.Fatal(err.Error())
	}
}

//userExsists checks if the username entered is already taken. If it is the function returns true.
func userExsist(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatal(err.Error())
	}
	count := 0

	for rows.Next() {
		count++
	}

	return count != 0
}
