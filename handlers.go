package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

// type Cookie struct {
// 	Name       string
// 	Value      string
// 	Path       string
// 	Domain     string
// 	Expires    time.Time
// 	RawExpires string
// 	// MaxAge=0 means no 'Max-Age' attribute specified.
// 	// MaxAgece means delete cookie now, equivalently 'Max-Age: 0'
// 	// MaxAge>e means Max-Age attribute present and given in seconds
// 	MaxAge   int
// 	Secure   bool
// 	HttpOnly bool
// 	// Samesite SameSite
// 	Raw      string
// 	Unparsed []string
// }

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

	if Person.Accesslevel {
		Executer(w, "templates/accessDenied.html")
	} else {
		tpl := template.Must(template.ParseGlob("templates/login.html"))
		if err := tpl.Execute(w, Person); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func LoginResult(w http.ResponseWriter, r *http.Request) {
	Person.Attempted = true
	email := r.FormValue("email")
	pass := r.FormValue("password")
	uuid := uuid.NewV4()

	if ValidEmail(email, sqliteDatabase) {
		if LoginValidator(email, pass, sqliteDatabase) {
			//Create the cookie

			if Person.Accesslevel {

				cookie, err := r.Cookie("1st-cookie")
				fmt.Println("cookie:", cookie, "err:", err)
				if err != nil {
					fmt.Println("cookie was not found")
					cookie = &http.Cookie{
						Name:     "1st-cookie",
						Value:    uuid.String(),
						HttpOnly: true,
						// MaxAge:   1000,
						Path: "/",
					}
					http.SetCookie(w, cookie)
					CookieAdd(cookie, sqliteDatabase)
				}
				// CookieAdd(cookie, sqliteDatabase)
			}
			Person.CookieChecker = true

			x := homePageStruct{MembersPost: Person, PostingDisplay: postData(sqliteDatabase)}
			tpl := template.Must(template.ParseGlob("templates/index.html"))
			if err := tpl.Execute(w, x); err != nil {
				log.Fatal(err.Error())
			}
		} else {
			tpl := template.Must(template.ParseGlob("templates/login.html"))
			if err := tpl.Execute(w, Person); err != nil {
				log.Fatal(err.Error())
			}
		}

	} else {
		tpl := template.Must(template.ParseGlob("templates/login2.html"))
		if err := tpl.Execute(w, Person); err != nil {
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

func Post(w http.ResponseWriter, r *http.Request) {
	Executer(w, "templates/newPost.html")

}

func postAdded(w http.ResponseWriter, r *http.Request) {
	cat := r.FormValue("categories")
	title := r.FormValue("title")
	post := r.FormValue("post")
	newPost(Person.Username, cat, title, post, sqliteDatabase)
	Executer(w, "templates/postAdded.html")

}

type homePageStruct struct {
	MembersPost    userDetails
	PostingDisplay []postDisplay
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Person)

	c1, err1 := r.Cookie("1st-cookie")

	if err1 == nil && !Person.Accesslevel {
		c1.MaxAge = -1
		http.SetCookie(w, c1)
	}

	c, err := r.Cookie("1st-cookie")

	if err != nil && Person.Accesslevel {
		//logged in and on 2nd browser
		Person.CookieChecker = false

	} else if err == nil && Person.Accesslevel {
		//Original browser
		Person.CookieChecker = true

	} else {
		// not logged in yet
		Person.CookieChecker = false
	}
	//Initialise the homePAgeStruct to pass through multiple data types
	x := homePageStruct{MembersPost: Person, PostingDisplay: postData(sqliteDatabase)}

	tpl := template.Must(template.ParseGlob("templates/index.html"))

	if err := tpl.Execute(w, x); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("YOUR COOKIE:", c)
}

func LogOut(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("1st-cookie")

	if err != nil {
		fmt.Println("Problem logging out with cookie")

	}

	c.MaxAge = -1
	http.SetCookie(w, c)

	if c.MaxAge == -1 {
		fmt.Println("Cookie deleted")
	}

	var newPerson userDetails
	Person = newPerson
	fmt.Println(Person)

	tpl := template.Must(template.ParseGlob("templates/logout.html"))

	if err := tpl.Execute(w, ""); err != nil {
		log.Fatal(err.Error())
	}
}
