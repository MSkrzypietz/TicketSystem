package webserver

import (
	"TicketSystem/config"
	"Ticketsystem/XML_IO"
	"fmt"
	"html/template"
	"net/http"
	//"github.com/stretchr/testify/assert"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromCookie(r)
	if RealUser(user) {
		// Show index Page
	} else {
		// Redirect to Login
		http.Redirect(w, r, "/login/", http.StatusFound)
	}
}

func StartServer() {
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/login/", ServeLogin)
	http.HandleFunc("/home/", ServeHome)
	http.HandleFunc("/logout/", ServeLogout)

	err := http.ListenAndServeTLS(":"+config.DEFAULT_PORT, config.CERT_FILE, config.KEY_FILE, nil)
	if err != nil {
		panic(err)
	}
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromCookie(r)
	if RealUser(user) {
		// Show home

	} else {
		// Redirect to login
		http.Redirect(w, r, "/login/", http.StatusFound)
	}
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromCookie(r)
	if !RealUser(user) {
		t, _ := template.ParseFiles("templates/login.html")
		fmt.Println(t.Execute(w, nil))

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")

		validUser := XML_IO.CheckUser(name, password)

		if validUser {
			fmt.Fprintf(w, "Hello, you're successfully logged in!")
			StartSession(w, name)
			http.Redirect(w, r, "/home/", http.StatusMovedPermanently)
		} else {
			fmt.Fprintf(w, "Something went wrong, please check your inputs")
			http.Redirect(w, r, "/login/", http.StatusMovedPermanently)
		}
		//
		//if err := scanner.Err(); err != nil {
		//	panic(err)
		//}
	} else {
		// User is already logged in
		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

func ServeLogout(w http.ResponseWriter, r *http.Request) {
	DestroySession(r)
	fmt.Fprintf(w, "You're logged out succesfully")
}
