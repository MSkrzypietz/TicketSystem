package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"TicketSystem/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

type context struct {
	Title           string
	ContentTemplate string
}

var ctx = context{Title: "Home", ContentTemplate: "home.html"}
var templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))

func StartServer() {
	// TODO: Fix paths and take it from configs..
	XML_IO.InitDataStorage("data/tickets", "data/users")

	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/signUp", ServeUserRegistration)
	http.HandleFunc("/signIn", ServeLogin)
	http.HandleFunc("/tickets/new", ServeNewTicket)
	http.HandleFunc("/createTicket", ServeTicketCreation)
	http.HandleFunc("/home", ServeHome)
	http.HandleFunc("/logout", ServeLogout)

	log.Printf("The server is starting to listen on port %d", config.Port)
	log.Printf("https://localhost:%d", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
	log.Println("The server has shutdown.")
}

func ServeNewTicket(w http.ResponseWriter, r *http.Request) {
	ctx = context{Title: "New Ticket", ContentTemplate: "newticket.html"}
	err := templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		log.Fatalf("Cannot Get View: %v", err)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	ctx = context{Title: "Home", ContentTemplate: "home.html"}
	err := templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		log.Fatalf("Cannot Get View: %v", err)
	}
}

func ServeTicketCreation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.PostFormValue("email") == "" {
		t, errTemplate := template.ParseFiles(path.Join(config.TemplatePath, "newticket.html"))
		if errTemplate != nil {
			t.Execute(w, nil)
		} else {
			log.Printf("Error accured with parsing html template newticket.html")
		}
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	email := r.PostFormValue("email")
	subject := r.PostFormValue("subject")
	message := r.PostFormValue("message")

	if utils.RegExMail(email) && utils.CheckEmptyXssString(subject) && utils.CheckEmptyXssString(message) {
		// Inputs okay

		// TODO: Handle errors from CreateTicket
		_, err = XML_IO.CreateTicket("data/tickets/ticket", "XML_IO/definitions.xml", email, subject, message)

		if err == nil {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			log.Printf("Error accured while creating a ticket: %d", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// Inputs not okay
		fmt.Fprintf(w, "Your Inputs are not valid. Please check your inputs and try again")
	}
}

func ServeUserRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if utils.CheckEqualStrings(r.PostFormValue("password1"), r.PostFormValue("password2")) {
		log.Println("Aborting registration... The entered passwords don't match.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password1")
	if utils.CheckUsernameFormal(username) && utils.CheckPasswdFormal(password) {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}

		_, errUser := XML_IO.CreateUser(config.UsersPath, username, string(hashedPassword))
		if errUser != nil {
			log.Println("Creating User failed, formal check of uname an passwd is valid")
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	} else {
		// Username or Password are formally not valid
		log.Println("Formal check of uname and passwd failed.")
		return
	}

}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	//user := GetUserFromCookie(r)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(path.Join(config.TemplatePath, "login.html"))
	t.Execute(w, nil)
}

func ServeLogout(w http.ResponseWriter, r *http.Request) {
	DestroySession(r)
	fmt.Fprintf(w, "You're logged out succesfully")
}
