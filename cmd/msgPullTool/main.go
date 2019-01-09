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
	// TODO: Replace with actual Get Request
	req := httptest.NewRequest(http.MethodGet, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webserver.ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	var mails utils.Response
	err := xml.NewDecoder(rr.Result().Body).Decode(&mails)
	if err != nil {
		log.Fatal(err)
	}

	if len(mails.Data) == 1 {
		fmt.Printf("There is %d E-Mail to be sent:\n", len(mails.Data))
	} else {
		fmt.Printf("There are %d E-Mails to be sent:\n", len(mails.Data))
	}

	for _, mail := range mails.Data {
		fmt.Println()
		fmt.Printf("E-Mail: %s\n", mail.EMailAddress)
		fmt.Printf("Subject: %s\n", mail.Subject)
		fmt.Printf("Message: %s\n", mail.Message)
	}
}
