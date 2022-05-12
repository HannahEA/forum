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

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// if Person.Accesslevel {
	// 	Executer(w, "templates/accessDenied.html")
	// } else {
	// 	tpl := template.Must(template.ParseGlob("templates/login.html"))
	// 	if err := tpl.Execute(w, Person); err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }

	tpl := template.Must(template.ParseGlob("templates/login.html"))
	if err := tpl.Execute(w, Person); err != nil {
		log.Fatal(err.Error())
	}
}

func LoginResult(w http.ResponseWriter, r *http.Request) {
	Person.Attempted = true
	email := r.FormValue("email")
	pass := r.FormValue("password")
	uuid := uuid.NewV4()

	if Person.Accesslevel {
		//The user is already logged in
		Person.Attempted = false
		tpl := template.Must(template.ParseGlob("templates/login.html"))
		if err := tpl.Execute(w, Person); err != nil {
			log.Fatal(err.Error())
		}
	} else if ValidEmail(email, sqliteDatabase) {
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
			Person.Attempted = false

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
		tpl := template.Must(template.ParseGlob("templates/login.html"))
		if err := tpl.Execute(w, Person); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/register.html"))
	if err := tpl.Execute(w, Person); err != nil {
		log.Fatal(err.Error())
	}
}

func registration2(w http.ResponseWriter, r *http.Request) {

	userN := r.FormValue("username")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	Person.RegistrationAttempted = true

	exist, _ := userExist(email, userN, sqliteDatabase)

	tpl := template.Must(template.ParseGlob("templates/register.html"))

	if exist {
		if err := tpl.Execute(w, Person); err != nil {
			log.Fatal(err.Error())
		}

	} else {
		Person.SuccessfulRegistration = true
		newUser(email, userN, pass, sqliteDatabase)

		if err := tpl.Execute(w, Person); err != nil {
			log.Fatal(err.Error())
		}

	}

}

func Post(w http.ResponseWriter, r *http.Request) {
	Person.PostAdded = false
	tpl := template.Must(template.ParseGlob("templates/newPost.html"))
	if err := tpl.Execute(w, Person); err != nil {
		log.Fatal(err.Error())
	}

}

func postAdded(w http.ResponseWriter, r *http.Request) {
	FEcat := r.FormValue("Frontend")
	BEcat := r.FormValue("BackEnd")
	FScat := r.FormValue("FullStack")

	cat := FEcat + " " + BEcat + " " + FScat
	//Loop through cat and remove any empty strings
	c := []rune(cat)
	category := []rune{}
	for i := 0; i < len(c); i++ {
		category = append(category, c[i])
		if c[i] == ' ' && c[i]+1 == ' ' {
			i++
		}
	}
	cat = string(category)
	title := r.FormValue("title")
	post := r.FormValue("post")
	newPost(Person.Username, cat, title, post, sqliteDatabase)
	//Add the post to the categories table with relevant table selected

	//Initialise the homePAgeStruct to pass through multiple data types
	x := homePageStruct{MembersPost: Person, PostingDisplay: postData(sqliteDatabase)}

	tpl := template.Must(template.ParseGlob("templates/index.html"))

	if err := tpl.Execute(w, x); err != nil {
		log.Fatal(err.Error())
	}

}

type homePageStruct struct {
	MembersPost    userDetails
	PostingDisplay []postDisplay
}

func Home(w http.ResponseWriter, r *http.Request) {
	//Likes
	postNum := r.FormValue("likeBtn")
	fmt.Printf("\n LIKE BUTTON VALUE \n")
	fmt.Println(postNum)
	LikeButton(postNum, sqliteDatabase)

	//Dislikes
	dislikePostNum := r.FormValue("dislikeBtn")
	fmt.Printf("\nDISLIKE BUTTON VALUE \n")
	DislikeButton(dislikePostNum, sqliteDatabase)

	// comments
	comment := r.FormValue("commentTxt")
	commentPostID := r.FormValue("commentSubmit")
	fmt.Printf("ADDING COMMENT: %v", commentPostID)
	newComment(Person.Username, commentPostID, comment, sqliteDatabase)

	//Comment likes
	commentNum := r.FormValue("commentlikeBtn")
	fmt.Printf("\n Comment LIKE BUTTON VALUE")
	fmt.Println(commentNum)
	CommentLikeButton(commentNum, sqliteDatabase)

	//Dislike comments
	commentDislike := r.FormValue("commentDislikeBtn")
	fmt.Printf("\nDISLIKE BUTTON VALUE \n")
	CommentDislikeButton(commentDislike, sqliteDatabase)

	//Make a button that gives a value depending on the filter button
	FE := r.FormValue("FEfilter")
	BE := r.FormValue("BEfilter")
	FS := r.FormValue("FSfilter")
	MyLikes := r.FormValue("likedPosts")
	Created := r.FormValue("myPosts")

	// all := r.FormValue("allfilter")
	postSlc := []postDisplay{}
	if FE == "FrontEnd" {
		frontEndSlc := []string{}
		//Create a query that gets the postIDs needed
		frontEndRows, errGetIDs := sqliteDatabase.Query("SELECT postID from categories WHERE FrontEnd = 1")
		if errGetIDs != nil {
			fmt.Println("EEROR trying to SELECT the posts with front end ID")
		}
		for frontEndRows.Next() {
			var GetIDs commentStruct

			err := frontEndRows.Scan(
				&GetIDs.CommentID,
			)
			if err != nil {
				fmt.Println("Error Scanning through rows")
			}

			frontEndSlc = append(frontEndSlc, GetIDs.CommentID)
		}
		postSlc = PostGetter(frontEndSlc, sqliteDatabase)

	} else if BE == "BackEnd" {
		BackEndSlc := []string{}
		//Create a query that gets the postIDs needed
		backEndRows, errGetIDs := sqliteDatabase.Query("SELECT postID from categories WHERE BackEnd = 1")
		if errGetIDs != nil {
			fmt.Println("EEROR trying to SELECT the posts with front end ID")
		}
		for backEndRows.Next() {
			var GetIDs commentStruct

			err := backEndRows.Scan(
				&GetIDs.CommentID,
			)
			if err != nil {
				fmt.Println("Error Scanning through rows")
			}

			BackEndSlc = append(BackEndSlc, GetIDs.CommentID)
		}
		postSlc = PostGetter(BackEndSlc, sqliteDatabase)

	} else if FS == "FullStack" {
		FullStackSlc := []string{}
		//Create a query that gets the postIDs needed
		FullStackRows, errGetIDs := sqliteDatabase.Query("SELECT postID from categories WHERE FullStack = 1")
		if errGetIDs != nil {
			fmt.Println("EEROR trying to SELECT the posts with front end ID")
		}
		for FullStackRows.Next() {
			var GetIDs commentStruct

			err := FullStackRows.Scan(
				&GetIDs.CommentID,
			)
			if err != nil {
				fmt.Println("Error Scanning through rows")
			}

			FullStackSlc = append(FullStackSlc, GetIDs.CommentID)
		}
		postSlc = PostGetter(FullStackSlc, sqliteDatabase)

	} else if MyLikes == "Liked Posts" {
		likedSlc := []string{}
		//Create a query that gets the postIDs needed
		likedRows, errGetIDs := sqliteDatabase.Query("SELECT postID from liketable WHERE reference = 1 AND user = (?)", Person.Username)
		if errGetIDs != nil {
			fmt.Println("EEROR trying to SELECT the posts with front end ID")
		}
		for likedRows.Next() {
			var GetIDs commentStruct

			err := likedRows.Scan(
				&GetIDs.CommentID,
			)
			if err != nil {
				fmt.Println("Error Scanning through rows")
			}

			likedSlc = append(likedSlc, GetIDs.CommentID)
		}
		postSlc = PostGetter(likedSlc, sqliteDatabase)
	} else if Created == "My Posts" {
		myPostsSlc := []string{}
		//Create a query that gets the postIDs needed
		myPostsRows, errGetIDs := sqliteDatabase.Query("SELECT postID from posts WHERE userName = (?)", Person.Username)
		if errGetIDs != nil {
			fmt.Println("EEROR trying to SELECT the posts with front end ID")
		}
		for myPostsRows.Next() {
			var GetIDs commentStruct

			err := myPostsRows.Scan(
				&GetIDs.CommentID,
			)
			if err != nil {
				fmt.Println("Error Scanning through rows")
			}

			myPostsSlc = append(myPostsSlc, GetIDs.CommentID)
		}
		postSlc = PostGetter(myPostsSlc, sqliteDatabase)

	}else {
		postSlc = postData(sqliteDatabase)
	}

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

	fmt.Printf("\n\n\n------------------------------------------------------------------Struct BEFORE: %v\n\n\n\n", Person)

	fmt.Printf("\n\n\n------------------------------------------------------------------Struct AFTER: %v\n\n\n\n", Person)
	//Initialise the homePAgeStruct to pass through multiple data types
	x := homePageStruct{MembersPost: Person, PostingDisplay: postSlc}

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
