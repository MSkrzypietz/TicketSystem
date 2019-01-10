package webserver

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"encoding/xml"
	"net/http"
	"path"
	"strconv"
	"time"
)

func ServeTickets(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil { // Show ticket overview
		ticketsData := utils.GetTicketsByStatus(utils.TicketStatusOpen)
		ticketsData = append(ticketsData, utils.GetTicketsByStatus(utils.TicketStatusInProcess)...)

		ctx := templateContext{HeaderTitle: "Tickets Overview", ContentTemplate: "tickets.html", IsSignedIn: true, Username: user.Username, TicketsData: ticketsData}
		executeTemplate(w, r, "index.html", ctx)
		return
	}

	ticket, err := utils.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorInvalidTicketID.ErrorPageURL(), http.StatusFound)
		return
	}

	// Creating a user list without the signed in user to show the selection for ticket assignment
	usersMap, err := utils.ReadUsers()
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageURL(), http.StatusFound)
		return
	}
	delete(usersMap, user.Username)
	usersList := []string{user.Username}
	for _, v := range usersMap {
		usersList = append(usersList, v.Username)
	}

	// TicketsData is used to display all possible tickets that can be merged, hence the current ticket gets removed
	ticketsData := utils.GetTicketsByEditor(user.Username)
	for i, t := range ticketsData {
		if t.Id == ticket.Id {
			ticketsData[i] = ticketsData[len(ticketsData)-1] // Replacing it with the last ticket
			ticketsData = ticketsData[:len(ticketsData)-1]   // Removing the last ticket
		}
	}
	ctx := templateContext{HeaderTitle: ticket.Reference, ContentTemplate: "ticketdetail.html", IsSignedIn: true, Username: user.Username, Users: usersList, TicketsData: ticketsData, CurrentTicket: ticket}
	executeTemplate(w, r, "index.html", ctx)
}

func ServeNewTicket(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := templateContext{HeaderTitle: "New Ticket", ContentTemplate: "newticket.html", IsSignedIn: err == nil}
	executeTemplate(w, r, "index.html", ctx)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	ctx := templateContext{HeaderTitle: "Home", ContentTemplate: "home.html", IsSignedIn: err == nil}
	executeTemplate(w, r, "index.html", ctx)
}

func ServeTicketCreation(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	email := r.PostFormValue("email")
	subject := r.PostFormValue("subject")
	message := r.PostFormValue("message")

	if email == "" || subject == "" || message == "" || !utils.RegExMail(email) || !utils.CheckEmptyXssString(subject) || !utils.CheckEmptyXssString(message) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageURL(), http.StatusFound)
		return
	}

	_, err := utils.CreateTicket(email, subject, message)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTicketCreation.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeUserRegistration(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	// Check if it has to show to template or if its a request to create one already
	if r.PostFormValue("username") == "" ||
		r.PostFormValue("password1") == "" ||
		r.PostFormValue("password2") == "" {
		ctx := templateContext{HeaderTitle: "Sign up", ContentTemplate: "signup.html", IsSignedIn: false}
		executeTemplate(w, r, "index.html", ctx)
		return
	}

	// Check if the passwords are not empty and if they are equal
	if !utils.CheckEqualStrings(r.PostFormValue("password1"), r.PostFormValue("password2")) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageURL(), http.StatusFound)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password1")

	// DebugMode removes annoying checks when testing
	if !config.DebugMode && (!utils.CheckUsernameFormal(username) || !utils.CheckPasswdFormal(password)) {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageURL(), http.StatusFound)
		return
	}

	_, err := utils.CreateUser(username, password)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUserCreation.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ServeAuthentication(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	// Check if it has to show to template or if its a request to sign in already
	if r.PostFormValue("username") == "" || r.PostFormValue("password") == "" {
		ctx := templateContext{HeaderTitle: "Sign in", ContentTemplate: "signin.html", IsSignedIn: false}
		executeTemplate(w, r, "index.html", ctx)
		return
	}

	uuid := utils.CreateUUID(64)
	// LoginUser checks if its the correct password for the username; if successful then save uuid as session cookie
	err := utils.LoginUser(r.PostFormValue("username"), r.PostFormValue("password"), uuid)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUserLogin.ErrorPageURL(), http.StatusFound)
		return
	}
	createSessionCookie(w, uuid)

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
	destroySession(w)
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
	if errCode > utils.ErrorCount()-1 || err != nil {
		http.Redirect(w, r, utils.ErrorUnknown.ErrorPageURL(), http.StatusMovedPermanently)
		return
	}

	ctx := templateContext{HeaderTitle: "Error", ContentTemplate: "errorpage.html", IsSignedIn: isSignedIn, ErrorMsg: utils.Error(errCode).String()}
	executeTemplate(w, r, "index.html", ctx)
}

func ServeAddComment(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	if len(r.PostFormValue("comment")) == 0 {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageURL(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	ticket, err := utils.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorInvalidTicketID.ErrorPageURL(), http.StatusFound)
		return
	}

	if r.PostFormValue("sendoption") == "comments" {
		_, err = utils.AddMessage(ticket, user.Username, r.PostFormValue("comment"))
		if err != nil {
			http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		}
	} else {
		err = utils.SendMail(ticket.Client, "Reply to your ticket (ID: "+strconv.Itoa(ticket.Id)+")", r.PostFormValue("comment"))
		if err != nil {
			http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)
}

func ServeTicketAssignment(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	// Check if the editor who is assigned to this ticket is an actual editor
	usersMap, err := utils.ReadUsers()
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageURL(), http.StatusFound)
		return
	}
	if _, ok := usersMap[r.PostFormValue("editor")]; !ok {
		http.Redirect(w, r, utils.ErrorInvalidInputs.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.ChangeEditor(ticketId, r.PostFormValue("editor"))
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.ChangeStatus(ticketId, utils.TicketStatusInProcess)
	if err != nil {
		// Removing the editor before showing the error page
		err = utils.ChangeEditor(ticketId, "")
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
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
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	ticket, err := utils.ReadTicket(ticketId)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataFetching.ErrorPageURL(), http.StatusFound)
		return
	}

	// Checking if the user is malicious and tries to release the ticket of someone else
	if ticket.Editor != user.Username {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.ChangeStatus(ticketId, utils.TicketStatusOpen)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.ChangeEditor(ticketId, "")
	if err != nil {
		// Reverting the status
		err = utils.ChangeStatus(ticketId, utils.TicketStatusInProcess)
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func ServeCloseTicket(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	ticketId, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.ChangeStatus(ticketId, utils.TicketStatusClosed)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/tickets/", http.StatusMovedPermanently)
}

func ServeMergeTickets(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	_, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	firstID, err := strconv.Atoi(path.Base(r.Referer()))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	secondID, err := strconv.Atoi(r.PostFormValue("ticket"))
	if err != nil {
		http.Redirect(w, r, utils.ErrorURLParsing.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.MergeTickets(firstID, secondID)
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)
}

func ServeChangeHolidayMode(w http.ResponseWriter, r *http.Request) {
	parseForm(w, r)

	user, err := utils.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, utils.ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		return
	}

	err = utils.SetUserHolidayMode(user.Username, r.PostFormValue("holidayMode") == "on")
	if err != nil {
		http.Redirect(w, r, utils.ErrorDataStoring.ErrorPageURL(), http.StatusFound)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)
}

func ServeMailsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// returns the list of mails which are to be sent
		getMails(w, r)
		return
	}

	if r.Method == http.MethodPost {
		// saves the posted mails to data storage
		postMails(w, r)
		return
	}

	utils.RespondWithError(w, http.StatusMethodNotAllowed, "This REST API only responds to GET and POST requests!")
}

func getMails(w http.ResponseWriter, _ *http.Request) {
	rawMails, err := utils.ReadMailsFile()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "We had issues fetching the E-Mails!")
		return
	}

	var mails []utils.MailData
	for _, mail := range rawMails.Maillist {
		if preventEMailPingPong(mail) {
			continue
		}

		err := mail.IncrementReadAttemptsCounter()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "We had internal issues fetching the data for you. Please try it again!")
			return
		} else {
			mail := utils.MailData{EMailAddress: mail.Mail, Subject: mail.Caption, Message: mail.Message}
			mails = append(mails, mail)
		}
	}

	utils.RespondWithXML(w, http.StatusOK, utils.Response{Meta: utils.MetaData{Code: http.StatusOK, Message: "OK"}, Data: mails})
}

// Checks depending on the FirstReadAttemptDate if the ReadAttemptCounter is tolerable to prevent email ping pong
func preventEMailPingPong(mail utils.Mail) bool {
	// There are no restraints for the first 3 read attempts
	if mail.ReadAttemptCounter < 3 {
		return false
	}

	// The next 7 read attempts can be made in the first 30 minutes after the first read attempt
	minsSinceFirstReadAttempt := time.Since(mail.FirstReadAttemptDate).Minutes()
	allowedAttempts := (float64(7)/30)*minsSinceFirstReadAttempt + 3
	if mail.ReadAttemptCounter < 10 {
		return allowedAttempts < float64(mail.ReadAttemptCounter)
	}

	// After that, every 30 minutes the mail can be read once more
	allowedAttempts = (float64(2)/30)*minsSinceFirstReadAttempt + 10
	return allowedAttempts < float64(mail.ReadAttemptCounter)
}

func postMails(w http.ResponseWriter, r *http.Request) {
	// Using MailData to ensure only accepting the address, subject and message
	var request utils.Request
	err := xml.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err = utils.CreateTicketFromMail(request.Mail.EMailAddress, request.Mail.Subject, request.Mail.Message)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "We had issues storing your sent E-Mails!")
		return
	}

	utils.RespondWithXML(w, http.StatusOK, utils.Response{Meta: utils.MetaData{Code: http.StatusOK, Message: "OK"}})
}

func ServeMailsSentNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "This REST API only responds to POST requests!")
		return
	}

	var request utils.Request
	err := xml.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	err = utils.DeleteMails(request.MailIDs)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "We ran into issues processing your request. Please try it again.")
		return
	}

	utils.RespondWithXML(w, http.StatusOK, utils.Response{Meta: utils.MetaData{Code: http.StatusOK, Message: "OK"}})
}

func executeTemplate(w http.ResponseWriter, r *http.Request, name string, ctx templateContext) {
	err := templates.ExecuteTemplate(w, name, ctx)
	if err != nil {
		http.Redirect(w, r, utils.ErrorTemplateExecution.ErrorPageURL(), http.StatusFound)
	}
}

func parseForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, utils.ErrorFormParsing.ErrorPageURL(), http.StatusFound)
	}
}
