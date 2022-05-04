package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/login.html"))
	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}

func LoginResult(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	if ValidEmail(email, sqliteDatabase) {
		if LoginValidator(email, pass, sqliteDatabase) {
			tpl := template.Must(template.ParseGlob("templates/homepage.html"))
			if err := tpl.Execute(w, ""); err != nil {
				log.Fatal(err.Error())
			}
		} else {
			tpl := template.Must(template.ParseGlob("templates/login2.html"))
			if err := tpl.Execute(w, "Incorrect password"); err != nil {
				log.Fatal(err.Error())
			}
		}

	} else {
		tpl := template.Must(template.ParseGlob("templates/login2.html"))
		if err := tpl.Execute(w, "Email not recognised"); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/register.html"))
	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/homepage.html"))
	p := Person
	fmt.Println(p)
	if err := tpl.Execute(w, p); err != nil {
		log.Fatal(err.Error())
	}
}
