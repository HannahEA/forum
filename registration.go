package main

import (
	"database/sql"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type userDetails struct {
	ID                     string
	Email                  string
	Username               string
	Password               string
	Accesslevel            bool
	CookieChecker          bool
	Attempted              bool
	RegistrationAttempted  bool
	FailedRegister         bool
	SuccessfulRegistration bool
	PostAdded              bool
}

//NewUser registers a new user to the database selected
func NewUser(email, username, password string, db *sql.DB) {
	hash, err := HashPassword(password)
	if err != nil {
		fmt.Println("Error hasing the password", err)
	}

	u1 := uuid.NewV4()
	_, errNewUser := db.Exec("INSERT INTO users (ID, email, username, password) VALUES (?, ?, ?, ?)", u1, email, username, hash)
	if errNewUser != nil {
		fmt.Printf("The error is %v", errNewUser.Error())
	}

}

//UserExsists checks if the username entered is already taken. If it is the function returns true.
func UserExist(email, username string, db *sql.DB) (bool, string) {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		fmt.Println("Error checking the email in userExist(): ", err)
	}
	count := 0

	for rows.Next() {
		count++
	}

	rows1, err1 := db.Query("SELECT username FROM users WHERE username = ?", username)
	if err1 != nil {
		fmt.Println("Error checking the username in userExist(): ", err1.Error())
	}

	count1 := 0
	for rows1.Next() {
		count1++
	}

	if count == 0 && count1 == 0 {
		return false, "You have successfully registered"
	} else if count1 == 1 && count == 1 {
		return true, "This account already exists. Please sign in"

	} else if count == 1 {
		return true, "This email is already taken. Please try a different email address"
	} else {
		return true, "This username is already taken. Please try a different username"
	}
}

//ValidEmail checks if the email entered exists in the database
func ValidEmail(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		fmt.Println("Error selecting the email in ValidEmail(): ", err)
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

//LoginValidaro checks if the email and passwords entered are the same
func LoginValidator(email, password string, db *sql.DB) bool {
	rows1, err1 := db.Query("SELECT ID, email, username, password FROM users WHERE email = ?", email)

	if err1 != nil {
		fmt.Println("Error selecting the details in LoginValidator(): ", err1.Error())
	}

	var u userDetails

	for rows1.Next() {
		err := rows1.Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.Password,
		)

		if err != nil {
			fmt.Println("SCANNING ERROR", err)
		}
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if hashErr == nil {
		Person.ID = u.ID
		Person.Email = u.Email
		Person.Username = u.Username
		Person.Password = u.Password
		Person.Accesslevel = true

	}

	return hashErr == nil

}

//HashPassword encrypts the password entered when registering
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CookieAdd(cookie *http.Cookie, db *sql.DB) {
	_, errCookie := db.Exec("INSERT INTO cookies (name, sessionID) VALUES (?, ?)", cookie.Name, cookie.Value)
	if errCookie != nil {
		fmt.Printf("The error is %v", errCookie.Error())
	}

}
