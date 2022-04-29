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
			"username" 	TEXT UNIQUE,
			"password"	TEXT 
		);
	`)

	if errTbl != nil {
		fmt.Println("TABLE ERROR")
		log.Fatal(errTbl.Error())
	}

	usersTbl.Exec()

	newUser("nater6", "HElloWorld", sqliteDatabase)

	fmt.Println(userExsist("nater68", sqliteDatabase))
	newUser("nater68", "HElloWorld", sqliteDatabase)

}

/*newUser adds a new account with a unique username and password to the database.

 */
func newUser(username, password string, db *sql.DB) {
	add, errNewUser := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if errNewUser != nil {
		log.Fatal()
	}

	_, preAdded := add.Exec(username, password)
	if preAdded != nil {
		log.Fatal(preAdded.Error())
	}
}

//userExsists checks if the username entered is already taken. If it is the function returns true.
func userExsist(username string, db *sql.DB) bool {
	rows, err := db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err.Error())
	}
	count := 0

	for rows.Next() {
		count++
	}

	return count != 0 
}
