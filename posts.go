package main

import (
	"database/sql"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
)

type postDisplay struct {
	PostID        string
	Username      string
	PostCategory  string
	Likes         int
	Dislikes      int
	TitleText     string
	PostText      string
	CookieChecker bool
}

//newPost creates a new post by a registered user
func newPost(userName, category, title, post string, db *sql.DB) {
	//If the title is empty the form is resubitting once all values have been reset so the post shouldn't be added to the database
	if title == "" {
		return
	}

	fmt.Println("ADDING POST")
	uuid := uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO posts (postID, userName, category, likes, dislikes, title, post) VALUES (?, ?, ?, 0, 0, ?, ?)", uuid, userName, category, title, post)
	if err != nil {
		log.Fatal(err.Error())
	}
	Person.PostAdded = true

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
		u.CookieChecker = Person.CookieChecker
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
		finalArray = append(finalArray, u)
	}
	return finalArray
}

func LikeButton(postID string, db *sql.DB) {
	//Check if the user has already liked this post/comment
	findRow, errRows := db.Query("SELECT reference FROM liketable WHERE postID = (?) AND user = (?)", postID, Person.Username)
	if errRows != nil {
		fmt.Println("SELECTING LIKE ERROR")
		log.Fatal(errRows.Error())
	}
	rounds := 0

	var check postDisplay
	for findRow.Next() {
		rounds++
		err2 := findRow.Scan(
			&check.Likes,
		)

		if err2 != nil {
			log.Fatal(err2.Error())
		}
	}

	//If rounds still equals 0 no row was found so we can insert the relevant row into our liketable
	if rounds == 0 {
		_, insertLikeErr := db.Exec("INSERT INTO liketable (user, postID, reference) VALUES (?, ?, 1)", Person.Username, postID)
		if insertLikeErr != nil {
			fmt.Println("Error when inserting into like table initially (LIKEBUTTON)")
			log.Fatal(insertLikeErr.Error())
		}

		//Increase likes
		LikeIncrease(postID, sqliteDatabase)
	} else {
		//Reference is equal to 1 so we need to undo the like action
		if check.Likes == 1 {
			LikeUndo(postID, sqliteDatabase)
			//Update reference to 0
			RefUpdate(0, postID, sqliteDatabase)
		} else if check.Likes == -1 {
			//user has already disliked so we must undislike the post and set it as liked
			DislikeUndo(postID, sqliteDatabase)
			LikeIncrease(postID, sqliteDatabase)
			//Update reference equal to 1

			RefUpdate(1, postID, sqliteDatabase)

		} else if check.Likes == 0 {
			//Increase likes only
			LikeIncrease(postID, sqliteDatabase)
			//set reference to 1
			RefUpdate(1, postID, sqliteDatabase)

		}
	}
}

func RefUpdate(value int, postID string, db *sql.DB) {
	_, err2 := db.Exec("UPDATE liketable SET reference = (?) WHERE postID = (?) AND user = (?)", value, postID, Person.Username)
	if err2 != nil {
		fmt.Println("UPDATING REFERENCE ")
		log.Fatal(err2.Error())
	}
}

func LikeIncrease(postID string, db *sql.DB) {
	//Increase likes

	likes, err := db.Query("SELECT likes FROM posts WHERE postID = (?)", postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	var temp postDisplay
	for likes.Next() {
		err := likes.Scan(
			&temp.Likes,
		)
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
	}

	temp.Likes++
	_, err2 := db.Exec("UPDATE posts SET likes = (?) WHERE postID = (?)", temp.Likes, postID)
	if err2 != nil {
		fmt.Println("UPDATING LIKES WHEN ROUNDS == 0")
		log.Fatal(err.Error())
	}

}

func LikeUndo(postID string, db *sql.DB) {
	likes, err := db.Query("SELECT likes FROM posts WHERE postID = (?)", postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	var temp postDisplay
	for likes.Next() {
		err := likes.Scan(
			&temp.Likes,
		)
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
	}

	temp.Likes--
	_, err2 := db.Exec("UPDATE posts SET likes = (?) WHERE postID = (?)", temp.Likes, postID)
	if err2 != nil {
		fmt.Println("LIKE UNDO")
		log.Fatal(err.Error())
	}
}

func DislikeIncrease(postID string, db *sql.DB) {
	dislikes, err := db.Query("SELECT dislikes FROM posts WHERE postID = (?)", postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	var temp postDisplay
	for dislikes.Next() {
		err := dislikes.Scan(
			&temp.Dislikes,
		)
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
	}

	temp.Dislikes++
	_, err2 := db.Exec("UPDATE posts SET dislikes = (?) WHERE postID = (?)", temp.Dislikes, postID)
	if err2 != nil {
		fmt.Println("UPDATING DISLIKES")
		log.Fatal(err.Error())
	}
}

func DislikeUndo(postID string, db *sql.DB) {
	dislikes, err := db.Query("SELECT dislikes FROM posts WHERE postID = (?)", postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	var temp postDisplay
	for dislikes.Next() {
		err := dislikes.Scan(
			&temp.Dislikes,
		)
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
	}

	temp.Dislikes--
	_, err2 := db.Exec("UPDATE posts SET dislikes = (?) WHERE postID = (?)", temp.Dislikes, postID)
	if err2 != nil {
		fmt.Println("DISLIKE UNDO")
		log.Fatal(err.Error())
	}
}

func DislikeButton(postID string, db *sql.DB) {
	//Check if the user has already liked/disliked this post/comment
	findRow, errRows := db.Query("SELECT reference FROM liketable WHERE postID = (?) AND user = (?)", postID, Person.Username)
	if errRows != nil {
		fmt.Println("SELECTING LIKE ERROR")
		log.Fatal(errRows.Error())
	}
	rounds := 0

	var check postDisplay
	for findRow.Next() {
		rounds++
		err := findRow.Scan(
			&check.Likes,
		)

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	//if rounds == 0 the user hasnt liked or disliked this post/comment yet
	if rounds == 0 {
		//Add the user to the liketable
		_, insertLikeErr := db.Exec("INSERT INTO liketable (user, postID, reference) VALUES (?, ?, -1)", Person.Username, postID)
		if insertLikeErr != nil {
			fmt.Println("Error when inserting into like table initially (DISLIKEBUTTON)")
			log.Fatal(insertLikeErr.Error())
		}
		//Increase number of dslikes
		DislikeIncrease(postID, sqliteDatabase)

	} else {
		if check.Likes == -1 {
			//The user has already disliked so we need to undo the dislike action
			DislikeUndo(postID, sqliteDatabase)
			//Change reference to 0
			RefUpdate(0, postID, sqliteDatabase)
		} else if check.Likes == 1 {
			//User has previously liked so we need to undo the like and dislike the comment
			//Undo like
			LikeUndo(postID, sqliteDatabase)
			// Increase dislike
			DislikeIncrease(postID, sqliteDatabase)
			//Set reference equal to -1
			RefUpdate(-1, postID, sqliteDatabase)
		} else if check.Likes == 0 {
			//The user is not currently liking or disliking the post so we need to increase dislike
			DislikeIncrease(postID, sqliteDatabase)
			//update reference to -1
			RefUpdate(-1, postID, sqliteDatabase)
		}
	}
}

// func DislikeButton2(postID string, db *sql.DB) {
// 	Dislikes, err := db.Query("SELECT dislikes FROM posts WHERE postID = (?)", postID)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	var temp postDisplay
// 	for Dislikes.Next() {
// 		err := Dislikes.Scan(
// 			&temp.Dislikes,
// 		)
// 		if err != nil {
// 			fmt.Println("SCANNING ERROR 342")
// 			log.Fatal(err.Error())
// 		}
// 	}
// 	temp.Dislikes++
// 	_, err2 := db.Exec("UPDATE posts SET dislikes = (?) WHERE postID = (?)", temp.Dislikes, postID)
// 	if err2 != nil {
// 		fmt.Println("DISLIKE ERROR")
// 		log.Fatal(err.Error())
// 	}

// }
