package main

import (
	"TicketSystem/utils"
	"TicketSystem/webserver"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	req := httptest.NewRequest(http.MethodGet, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webserver.ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	var mails utils.Response
	err := xml.NewDecoder(rr.Result().Body).Decode(&mails)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("There are %d E-Mails to be sent:\n", len(mails.Data))
	for _, mail := range mails.Data {
		fmt.Println()
		fmt.Printf("ID: %d\n", mail.MailId)
		fmt.Printf("E-Mail: %s\n", mail.Mail)
		fmt.Printf("Subject: %s\n", mail.Caption)
		fmt.Printf("Message: %s\n", mail.Message)
	}
}
