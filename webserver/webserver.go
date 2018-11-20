package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	//"github.com/stretchr/testify/assert"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	//user := GetUserFromCookie(r)
	//if RealUser(user) {
	//	// Show index Page
	//} else {
	// Redirect to Login
	http.Redirect(w, r, "/login", http.StatusFound)
	//}
}

func StartServer() {
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/register", ServeUserRegistration)
	http.HandleFunc("/login", ServeLogin)
	http.HandleFunc("/home", ServeHome)
	http.HandleFunc("/logout", ServeLogout)

	err := http.ListenAndServeTLS(":"+config.Port, config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
}

func ServeUserRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Refactored into own func
	if r.PostFormValue("newPassword1") != r.PostFormValue("newPassword2") {
		log.Println("Aborting registration... The entered passwords don't match.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: More validation check like username == pw? or len(username) > 4? ...

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("newPassword1")), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// TODO: Handle error from CreateUser
	XML_IO.CreateUser(config.UsersPath, r.PostFormValue("newUsername"), string(hashedPassword))
	w.WriteHeader(http.StatusAccepted) // if no error

	http.Redirect(w, r, "/", http.StatusFound)
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromCookie(r)
	if RealUser(user) {
		// Show home

	} else {
		// Redirect to login
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	//user := GetUserFromCookie(r)
	//if !RealUser(user) {
	t, _ := template.ParseFiles("templates/login.html")
	fmt.Println(t.Execute(w, nil))
	//
	//	err := r.ParseForm()
	//	if err != nil {
	//		panic(err)
	//	}
	//	name := r.PostFormValue("name")
	//	password := r.PostFormValue("password")
	//
	//	validUser := CheckUser(name, password)
	//
	//	if validUser {
	//		fmt.Fprintf(w, "Hello, you're successfully logged in!")
	//		StartSession(w, name)
	//		http.Redirect(w, r, "/home/", http.StatusMovedPermanently)
	//	} else {
	//		fmt.Fprintf(w, "Something went wrong, please check your inputs")
	//		http.Redirect(w, r, "/login/", http.StatusMovedPermanently)
	//	}
	//	//
	//	//if err := scanner.Err(); err != nil {
	//	//	panic(err)
	//	//}
	//} else {
	//	// User is already logged in
	//	http.Redirect(w, r, "/home/", http.StatusFound)
	//}
}

func ServeLogout(w http.ResponseWriter, r *http.Request) {
	DestroySession(r)
	fmt.Fprintf(w, "You're logged out succesfully")
}
