package main

import (
	"database/sql"
	"fmt"
	"log"
)

type postDisplay struct {
	postID       int
	username     string
	postCategory string
	likes        int
	dislikes     int
	titleText    string
	postText     string
}

//newPost creates a new post by a registered user
func newPost(userName, category, title, post string, db *sql.DB) {
	_, err := db.Exec("INSERT INTO posts (userName, category, likes, dislikes, title, post) VALUES (?,?, 0,0, ?, ?)", userName, category, title, post)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func postData(db *sql.DB) []postDisplay {
	rows, err := db.Query("SELECT postID, userName, category, likes, dislikes, title, post FROM posts")
	if err != nil {
		log.Fatal(err.Error())
	}

	finalArray := []postDisplay{}

	for rows.Next() {

		var u postDisplay
		err := rows.Scan(
			&u.postID,
			&u.username,
			&u.postCategory,
			&u.likes,
			&u.dislikes,
			&u.titleText,
			&u.postText,
		)

		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
		finalArray = append(finalArray, u)
	}
	return finalArray
}
