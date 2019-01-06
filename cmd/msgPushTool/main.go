package main

import (
	"TicketSystem/XML_IO"
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

	mail := XML_IO.Mail{Mail: *mailAddress, Caption: *subject, Message: *msg}
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

	if response.Success {
		fmt.Println("Successfully received the email.")
	} else {
		fmt.Println("We were not able to store your email. Please try it again!")
	}
}
