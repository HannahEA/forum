package main

import (
	"database/sql"
	"log"
)

//newPost creates a new post by a registered user
func newPost(userName, category, title, post string, db *sql.DB) {
	_, err := db.Exec("INSERT INTO posts (userName, category, likes, dislikes, title, post) VALUES (?,?, 0,0, ?, ?)",userName, category, title, post)
	if err != nil {
		log.Fatal(err.Error())
	}
}

