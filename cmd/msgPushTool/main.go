package main

import (
	"TicketSystem/utils"
	"TicketSystem/webserver"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	mailAddress := flag.String("mail", "test@gmail.com", "Email Address")
	subject := flag.String("subject", "Test Caption", "Ticket Subject")
	msg := flag.String("message", "Test Message", "Ticket Message")
	flag.Parse()

	mail := utils.Mail{Mail: *mailAddress, Caption: *subject, Message: *msg}
	buf, err := xml.Marshal(mail)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/mails", bytes.NewBuffer(buf))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webserver.ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	var response utils.Response
	err = xml.NewDecoder(rr.Result().Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	if response.Meta.Code == http.StatusOK {
		fmt.Println("Successfully received the email with the following content:")
		fmt.Printf("\tE-Mail: %s\n", mail.Mail)
		fmt.Printf("\tSubject: %s\n", mail.Caption)
		fmt.Printf("\tMessage: %s\n", mail.Message)
	} else {
		fmt.Println("We were not able to store your email. Please try it again!")
	}
}
