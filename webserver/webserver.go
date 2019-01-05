package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"TicketSystem/utils"
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
	Username        string
	Users           []string
	TicketsData     []XML_IO.Ticket
}

var templates *template.Template

func Setup() {
	err := XML_IO.InitDataStorage()
	if err != nil {
		log.Fatal("Cannot start the ticket system due to problems initializing the data storage...")
	}
	templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))
}

func StartServer() {
	Setup()
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
	http.HandleFunc("/error/", ServeErrorPage)
	http.HandleFunc("/addComment", ServeAddComment)
	http.HandleFunc("/assignTicket", ServeTicketAssignment)
	http.HandleFunc("/releaseTicket", ServeTicketRelease)

	log.Printf("The server is starting to listen on https://localhost:%d", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
}

func ServeTickets(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil { // Show ticket overview
		// First listing all the tickets the current user is assigned to and afterwards showing not assigned tickets
		ticketsData := XML_IO.GetTicketsByEditor(user.Username)
		ticketsData = append(ticketsData, XML_IO.GetTicketsByStatus(0)...)

		ctx := context{HeaderTitle: "Tickets Overview", ContentTemplate: "tickets.html", IsSignedIn: true, Username: user.Username, TicketsData: ticketsData}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		if err != nil {
			http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
		}
		return
	}

	ticket, err := XML_IO.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorInvalidTicketID.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Checking if this ticket is assigned to the current user
	if ticket.Status != XML_IO.UnProcessed && ticket.Editor != user.Username {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	usersMap, err := XML_IO.ReadUsers()
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageUrl(), http.StatusFound)
		return
	}
	delete(usersMap, user.Username)

	usersList := []string{user.Username}
	for _, v := range usersMap {
		usersList = append(usersList, v.Username)
	}

	ctx := context{HeaderTitle: ticket.Reference, ContentTemplate: "ticketdetail.html", IsSignedIn: true, Username: user.Username, Users: usersList, TicketsData: []XML_IO.Ticket{ticket}}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
	}
}

func ServeNewTicket(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := context{HeaderTitle: "New Ticket", ContentTemplate: "newticket.html", IsSignedIn: err == nil}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := context{HeaderTitle: "Home", ContentTemplate: "home.html", IsSignedIn: err == nil}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
	}
}

func ServeTicketCreation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	email := r.PostFormValue("email")
	subject := r.PostFormValue("subject")
	message := r.PostFormValue("message")

	if email == "" || subject == "" || message == "" || !utils.RegExMail(email) || !utils.CheckEmptyXssString(subject) || !utils.CheckEmptyXssString(message) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	_, err = XML_IO.CreateTicket(email, subject, message)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTicketCreation.ErrorPageUrl(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeUserRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Check if it has to show to template or if its a request to create one already
	if r.PostFormValue("username") == "" ||
		r.PostFormValue("password1") == "" ||
		r.PostFormValue("password2") == "" {
		ctx := context{HeaderTitle: "Sign up", ContentTemplate: "signup.html", IsSignedIn: false}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		if err != nil {
			http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
		}
		return
	}

	// Check if the passwords are not empty and if they equal
	if !utils.CheckEqualStrings(r.PostFormValue("password1"), r.PostFormValue("password2")) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password1")
	if utils.CheckUsernameFormal(username) && utils.CheckPasswdFormal(password) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	_, err = XML_IO.CreateUser(username, password)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUserCreation.ErrorPageUrl(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Check if it has to show to template or if its a request to sign in already
	if r.PostFormValue("username") == "" || r.PostFormValue("password") == "" {
		ctx := context{HeaderTitle: "Sign in", ContentTemplate: "signin.html", IsSignedIn: false}
		err = templates.ExecuteTemplate(w, "index.html", ctx)
		if err != nil {
			http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
		}
		return
	}

	uuid := utils.CreateUUID(64)
	err = XML_IO.LoginUser(r.PostFormValue("username"), r.PostFormValue("password"), uuid)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUserLogin.ErrorPageUrl(), http.StatusFound)
		return
	}
	CreateSessionCookie(w, uuid)

	// This will redirect the user to his original destination if he was forced to authorize
	url, err := r.Cookie("requested-url-while-not-authenticated")
	if err != nil || url.Value == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, url.Value, http.StatusFound)
	}
}

func ServeSignOut(w http.ResponseWriter, r *http.Request) {
	// Destroying user specific cookies
	DestroySession(w)
	http.SetCookie(w, &http.Cookie{
		Name:   "requested-url-while-not-authenticated",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeErrorPage(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	isSignedIn := err == nil

	// Error can be ignored, because the errCode will be set to 0 and hence a correct error page will be displayed
	errCode, _ := strconv.Atoi(path.Base(r.URL.Path))

	ctx := context{HeaderTitle: "Error", ContentTemplate: "errorpage.html", IsSignedIn: isSignedIn, ErrorMsg: utils.Error(errCode).String()}
	err = templates.ExecuteTemplate(w, "index.html", ctx)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageUrl(), http.StatusFound)
	}
}

func ServeAddComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticket, err := XML_IO.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageUrl(), http.StatusFound)
		return
	}

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageUrl(), http.StatusFound)
		return
	}

	_, err = XML_IO.AddMessage(ticket, user.Username, r.PostFormValue("comment"))
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
	}

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func ServeTicketAssignment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Check if the editor who is assigned to this ticket is an actual editor
	usersMap, err := XML_IO.ReadUsers()
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageUrl(), http.StatusFound)
		return
	}
	if _, ok := usersMap[r.PostFormValue("editor")]; !ok {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	err = XML_IO.ChangeEditor(ticketId, r.PostFormValue("editor"))
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		return
	}

	err = XML_IO.ChangeStatus(ticketId, XML_IO.InProcess)
	if err != nil {
		// Removing the editor before showing the error page
		err = XML_IO.ChangeEditor(ticketId, "")
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		return
	}

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Resides on the ticket when assigned to oneself, else the user gets send to the tickets overview
	if r.PostFormValue("editor") == user.Username {
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	} else {
		http.Redirect(w, r, "/tickets/", http.StatusFound)
	}
}

func ServeTicketRelease(w http.ResponseWriter, r *http.Request) {
	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticket, err := XML_IO.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageUrl(), http.StatusFound)
		return
	}

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	// Checking if the user is malicious
	if ticket.Editor != user.Username {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	err = XML_IO.ChangeStatus(ticketId, XML_IO.UnProcessed)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		return
	}

	err = XML_IO.ChangeEditor(ticketId, "")
	if err != nil {
		// Reverting the status
		err = XML_IO.ChangeStatus(ticketId, XML_IO.InProcess)
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
