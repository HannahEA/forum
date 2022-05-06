package main

import (
	"database/sql"
	"fmt"
	"log"
)

type postDisplay struct {
	PostID       int
	Username     string
	PostCategory string
	Likes        int
	Dislikes     int
	TitleText    string
	PostText     string
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
			&u.PostID,
			&u.Username,
			&u.PostCategory,
			&u.Likes,
			&u.Dislikes,
			&u.TitleText,
			&u.PostText,
		)

		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
		finalArray = append(finalArray, u)
	}
	return finalArray
}

func LikeButton(postID string, db *sql.DB) {
	likes, err:= db.Query("SELECT Likes FROM posts WHERE postID = (?)", postID)
	if err != nil {
		log.Fatal(err.Error())
	}
	var temp postDisplay
	for likes.Next() {
		err := likes.Scan (
			&temp.Likes,
		)
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
	}
	temp.Likes++
	_, err2:= db.Exec("UPDATE posts SET likes = (?) WHERE postID = (?)", temp.Likes, postID)
	if err2 != nil {
		fmt.Println("LIKE ERROR")
		log.Fatal(err.Error())
	}
	
	
}
