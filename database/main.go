package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
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
	_, errTbl := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	TEXT,
			"email" 	TEXT UNIQUE,
			"username"	TEXT UNIQUE,
			"password"	TEXT 
		);
	`)

	if errTbl != nil {
		fmt.Println("USER ERROR")
		log.Fatal(errTbl.Error())
	}

	// usersTbl.Exec()

	//Create the posts table
	_, errPosts := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"postID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"userName"	TEXT NOT NULL,
			"category"	TEXT ,
			"likes" INTEGER,
			"dislikes" INTEGER,
			"title" TEXT,
			"post" TEXT
		);
	`)

	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}

	//Create the cookies table
	_, errCookie := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "cookies" (
			"name"	TEXT,
			"sessionID" 	TEXT UNIQUE
		);
	`)

	if errCookie != nil {
		fmt.Println("TABLE ERROR")
		log.Fatal(errTbl.Error())
	}
}
