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

//TODO: Global context won t work when multiple users are requesting stuff
// Probably need to create a SessionManager struct and perhaps a router struct
// Remove context all together, instead have different name when executing a template

type context struct {
	HeaderTitle     string
	ContentTemplate string
	IsSignedIn      bool
	ShowSignInModal bool
}

//var ctx context// = context{HeaderTitle: "Home", ContentTemplate: "home.html", IsSignedIn: false, ShowSignInModal: true}
var templates *template.Template

func StartServer() {
	// TODO: Fix paths and take it from configs..
	XML_IO.InitDataStorage("data/tickets", "data/users")
	templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))

	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/signUp", ServeUserRegistration)
	http.HandleFunc("/signIn", ServeSignIn)
	http.HandleFunc("/signOut", ServeSignOut)
	http.HandleFunc("/tickets/", ServeTickets)
	http.HandleFunc("/tickets/new", ServeNewTicket)
	http.HandleFunc("/createTicket", ServeTicketCreation)

	log.Printf("The server is starting to listen on https://localhost:%d", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
	log.Println("The server has shutdown.")
}

func ServeTickets(w http.ResponseWriter, r *http.Request) {
	//TODO: Replace (the following 5 lines) with wrapper that checks if its a valid session (to access restrained resources
	_, err := GetUserFromCookie(r)
	if err != nil {
		ctx := context{HeaderTitle: "Tickets Overview", ContentTemplate: "tickets.html", IsSignedIn: false, ShowSignInModal: true}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		return
	}

	ctx := context{HeaderTitle: "Tickets Overview", ContentTemplate: "tickets.html", IsSignedIn: true, ShowSignInModal: false}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		log.Fatalf("Cannot Get View: %v", err)
	}
}

func ServeNewTicket(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserFromCookie(r)
	ctx := context{HeaderTitle: "New Ticket", ContentTemplate: "newticket.html", IsSignedIn: err == nil, ShowSignInModal: false}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		log.Fatalf("Cannot Get View: %v", err)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserFromCookie(r)
	ctx := context{HeaderTitle: "Home", ContentTemplate: "home.html", IsSignedIn: err == nil, ShowSignInModal: false}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
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

	if !utils.CheckEqualStrings(r.PostFormValue("password1"), r.PostFormValue("password2")) {
		log.Println("Aborting registration... The entered passwords don't match.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password1")
	//TODO: uncomment in production
	// if utils.CheckUsernameFormal(username) && utils.CheckPasswdFormal(password) {
	if true {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}

		_, errUser := XML_IO.CreateUser(config.UsersFilePath, username, string(hashedPassword))
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

func ServeSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := utils.CreateUUID(64) // TODO: put this in LoginUser
	err = XML_IO.LoginUser(config.UsersFilePath, r.PostFormValue("username"), r.PostFormValue("password"), uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	CreateCookie(w, uuid) // TODO: put this in LoginUser

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func ServeSignOut(w http.ResponseWriter, r *http.Request) {
	DestroySession(w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
