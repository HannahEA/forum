package main

import (
	"database/sql"
	"log"
)

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

