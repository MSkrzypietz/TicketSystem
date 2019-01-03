package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"TicketSystem/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

type context struct {
	HeaderTitle     string
	ContentTemplate string
	IsSignedIn      bool
	ErrorMsg        string
}

var templates *template.Template

func StartServer() {
	XML_IO.InitDataStorage()
	templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))

	af := utils.AuthenticatorFunc(func(user, password string) bool {
		ok, err := XML_IO.VerifyUser(user, password)
		if err != nil {
			log.Println(err)
		}
		return ok
	})

	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/signUp", ServeUserRegistration)
	http.HandleFunc("/signIn", ServeSignIn)
	http.HandleFunc("/signOut", ServeSignOut)
	http.HandleFunc("/tickets/", utils.AuthWrapper(af, ServeTickets))
	http.HandleFunc("/tickets/new", ServeNewTicket)
	http.HandleFunc("/createTicket", ServeTicketCreation)
	http.HandleFunc("/error", ServeErrorPage)

	log.Printf("The server is starting to listen on https://localhost:%d", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
	log.Println("The server has shutdown.")
}

func ServeTickets(w http.ResponseWriter, r *http.Request) {
	ctx := context{HeaderTitle: "Tickets Overview", ContentTemplate: "tickets.html", IsSignedIn: true}
	err := templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		// TODO: How to handle? Fatal should be avoided
	}
}

func ServeNewTicket(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := context{HeaderTitle: "New Ticket", ContentTemplate: "newticket.html", IsSignedIn: err == nil}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		log.Fatalf("Cannot Get View: %v", err)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := context{HeaderTitle: "Home", ContentTemplate: "home.html", IsSignedIn: err == nil}
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
		_, err = XML_IO.CreateTicket(email, subject, message)

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

	if r.PostFormValue("username") == "" ||
		r.PostFormValue("password1") == "" ||
		r.PostFormValue("password2") == "" {
		ctx := context{HeaderTitle: "Sign up", ContentTemplate: "signup.html", IsSignedIn: false}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		if err != nil {
			//TODO
		}
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
		_, errUser := XML_IO.CreateUser(username, string(password))
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

	if r.PostFormValue("username") == "" || r.PostFormValue("password") == "" {
		ctx := context{HeaderTitle: "Sign in", ContentTemplate: "signin.html", IsSignedIn: false}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		if err != nil {
			//TODO
		}
		return
	}

	uuid := utils.CreateUUID(64) // TODO: put this in LoginUser
	err = XML_IO.LoginUser(r.PostFormValue("username"), r.PostFormValue("password"), uuid)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error?err=Test", http.StatusFound)
		return
	}
	CreateCookie(w, uuid) // TODO: put this in LoginUser

	url, err := r.Cookie("requested-url-while-not-authenticated")
	if err != nil || url.Value == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, url.Value, http.StatusFound)
	}
}

func ServeSignOut(w http.ResponseWriter, r *http.Request) {
	DestroySession(w)
	http.SetCookie(w, &http.Cookie{
		Name:   "requested-url-while-not-authenticated",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeErrorPage(w http.ResponseWriter, r *http.Request) {
	// TODO: Fix IsSignedIn
	ctx := context{HeaderTitle: "Error", ContentTemplate: "errorpage.html", IsSignedIn: false, ErrorMsg: r.URL.Query()["err"][0]}
	err := templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {

	}
}
