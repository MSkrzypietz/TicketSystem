package webserver

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"context"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

type templateContext struct {
	HeaderTitle     string
	ContentTemplate string
	IsSignedIn      bool
	IsUserInHoliday bool
	ErrorMsg        string
	Username        string
	Users           []utils.User
	TicketsData     []utils.Ticket
	CurrentTicket   utils.Ticket
}

var templates *template.Template

func Setup() {
	err := utils.InitDataStorage()
	if err != nil {
		log.Fatal("Cannot start the ticket system due to problems initializing the data storage...")
	}
	templates = template.Must(template.ParseGlob(path.Join(config.TemplatePath, "*")))
}

func StartServer(done <-chan bool, shutdown chan<- bool) {
	Setup()

	// Using http.NewServeMux() to prevent panics for multiple registrations when testing the cli tools
	handler := http.NewServeMux()
	handler.HandleFunc("/", ServeIndex)
	handler.HandleFunc("/signUp", ServeUserRegistration)
	handler.HandleFunc("/signIn", ServeAuthentication)
	handler.HandleFunc("/signOut", ServeSignOut)
	handler.HandleFunc("/tickets/", authenticate(ServeTickets))
	handler.HandleFunc("/tickets/new", ServeNewTicket)
	handler.HandleFunc("/createTicket", ServeTicketCreation)
	handler.HandleFunc("/error/", ServeErrorPage)
	handler.HandleFunc("/addComment", authenticate(ServeAddComment))
	handler.HandleFunc("/assignTicket", authenticate(ServeTicketAssignment))
	handler.HandleFunc("/releaseTicket", authenticate(ServeTicketRelease))
	handler.HandleFunc("/closeTicket", authenticate(ServeCloseTicket))
	handler.HandleFunc("/mergeTickets", authenticate(ServeMergeTickets))
	handler.HandleFunc("/changeHolidayMode", authenticate(ServeChangeHolidayMode))
	handler.HandleFunc("/mails", ServeMailsAPI)
	handler.HandleFunc("/mails/notify", ServeMailsSentNotification)

	server := &http.Server{Addr: "localhost:" + strconv.Itoa(config.Port), Handler: handler}

	go func() {
		log.Printf("The server is starting to listen on https://localhost:%d", config.Port)
		err := server.ListenAndServeTLS(config.ServerCertPath, config.ServerKeyPath)
		if err != nil {
			log.Println(err)
		}
	}()

	// Waiting for user input to begin shutting down the server
	<-done

	log.Println("Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down the server: %v\n", err)
	}
	log.Println("The shut down gracefully :)")

	// Sending a signal back to let the server shut down gracefully and not get interrupted by the main function existing
	shutdown <- true
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
