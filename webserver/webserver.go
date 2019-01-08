package webserver

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

type templateContext struct {
	HeaderTitle     string
	ContentTemplate string
	IsSignedIn      bool
	ErrorMsg        string
	Username        string
	Users           []string
	TicketsData     []utils.Ticket
}

var templates *template.Template

// TODO: Try replacing this with init()
func Setup() {
	err := utils.InitDataStorage()
	if err != nil {
		log.Fatal("Cannot start the ticket system due to problems initializing the data storage...")
	}
	templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))
}

func StartServer() {
	Setup()

	// TODO: Only allow certain http methods create a wrapper in this file
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/signUp", ServeUserRegistration)
	http.HandleFunc("/signIn", ServeAuthentication)
	http.HandleFunc("/signOut", ServeSignOut)
	http.HandleFunc("/tickets/", authenticate(ServeTickets))
	http.HandleFunc("/tickets/new", ServeNewTicket)
	http.HandleFunc("/createTicket", ServeTicketCreation)
	http.HandleFunc("/error/", ServeErrorPage)
	http.HandleFunc("/addComment", authenticate(ServeAddComment))
	http.HandleFunc("/assignTicket", authenticate(ServeTicketAssignment))
	http.HandleFunc("/releaseTicket", authenticate(ServeTicketRelease))
	http.HandleFunc("/closeTicket", authenticate(ServeCloseTicket))
	http.HandleFunc("/mails", ServeMailsAPI)
	http.HandleFunc("/mails/notify", ServeMailsSentNotification)

	log.Printf("The server is starting to listen on https://localhost:%d", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.ServerCertPath, config.ServerKeyPath, nil)
	if err != nil {
		panic(err)
	}
}

// Wrapper to check for session cookie
func authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserFromCookie(r)
		// Checks if session cookie is available and if its a correct one
		if err != nil || utils.VerifySessionCookie(user.Username, user.Password) != nil {
			// Setting a cookie to remember the requested url from the user
			// in order to redirect him there after successful authentication
			utils.RemoveCookie(w, "requested-url-while-not-authenticated")
			http.SetCookie(w, &http.Cookie{
				Name:     "requested-url-while-not-authenticated",
				Value:    r.URL.RequestURI(),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60,
			})

			http.Redirect(w, r, "/signIn", http.StatusFound)
			return
		}

		handler(w, r)
	}
}
