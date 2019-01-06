package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"TicketSystem/utils"
	"encoding/xml"
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
	http.HandleFunc("/addComment", utils.AuthWrapper(af, ServeAddComment))
	http.HandleFunc("/assignTicket", utils.AuthWrapper(af, ServeTicketAssignment))
	http.HandleFunc("/releaseTicket", utils.AuthWrapper(af, ServeTicketRelease))
	http.HandleFunc("/mails", ServeMailsAPI)
	http.HandleFunc("/mails/notify", ServeMailsSentNotification)

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
		ticketsData := XML_IO.GetTicketsByStatus(XML_IO.UnProcessed)
		ticketsData = append(ticketsData, XML_IO.GetTicketsByStatus(XML_IO.InProcess)...)

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

	// Creating a user list without the signed in user to show the selection for ticket assignment
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

	// Check if the passwords are not empty and if they are equal
	if !utils.CheckEqualStrings(r.PostFormValue("password1"), r.PostFormValue("password2")) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password1")

	// @Neemann: In Produktion würden wird das auskommentieren, aber zum testen wäre es sonst zu nervig..
	//if !utils.CheckUsernameFormal(username) || !utils.CheckPasswdFormal(password) {
	if false {
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

	errCode, err := strconv.Atoi(path.Base(r.URL.Path))
	if errCode > utils.ErrorCount-1 || err != nil {
		http.Redirect(w, r, utils.ErrorUnknown.ErrorPageUrl(), http.StatusMovedPermanently)
		return
	}

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

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
		return
	}

	if len(r.PostFormValue("comment")) == 0 {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	ticket, err := XML_IO.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorInvalidTicketID.ErrorPageUrl(), http.StatusFound)
		return
	}

	if r.PostFormValue("sendoption") == "comments" {
		_, err = XML_IO.AddMessage(ticket, user.Username, r.PostFormValue("comment"))
		if err != nil {
			http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		}
	} else {
		err = XML_IO.SendMail(ticket.Client, "Reply to your ticket (ID: "+strconv.Itoa(ticket.Id)+")", r.PostFormValue("comment"))
		if err != nil {
			http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageUrl(), http.StatusFound)
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)
}

func ServeTicketAssignment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageUrl(), http.StatusFound)
		return
	}

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
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

	// Resides on the ticket when assigned to oneself, else the user gets send to the tickets overview
	if r.PostFormValue("editor") == user.Username {
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	} else {
		http.Redirect(w, r, "/tickets/", http.StatusFound)
	}
}

func ServeTicketRelease(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageUrl(), http.StatusFound)
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

	// Checking if the user is malicious and tries to release the ticket of someone else
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

func ServeMailsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// returns the list of mails which are to be sent
		getMails(w, r)
		return
	}

	utils.RespondWithError(w, http.StatusMethodNotAllowed, "This REST API only responds to GET and POST requests!")
}

func getMails(w http.ResponseWriter, _ *http.Request) {
	rawMails, err := XML_IO.ReadMailsFile()
	if err != nil {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "We had issues fetching the E-Response!")
		return
	}

	utils.RespondWithXML(w, http.StatusOK, utils.Response{Success: true, Data: rawMails.Maillist})
}

func ServeMailsSentNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "This REST API only responds to POST requests!")
	}

	var mails utils.MailIDs
	err := xml.NewDecoder(r.Body).Decode(&mails)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
	}

	err = XML_IO.DeleteMails(mails.IDList)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "We ran into issues processing your request. Please try it again.")
	}

	utils.RespondWithXML(w, http.StatusOK, utils.Response{Success: true})
}
