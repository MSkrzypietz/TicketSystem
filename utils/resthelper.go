package utils

import (
	"encoding/xml"
	"log"
	"net/http"
)

// this file is inspired by https://itnext.io/building-restful-web-api-service-using-golang-chi-mysql-d85f427dee54

type Request struct {
	Mail    MailData `xml:"mail,omitempty"`
	MailIDs []int    `xml:"mails>mailID,omitempty"`
}

type Response struct {
	Meta MetaData   `xml:"meta"`
	Data []MailData `xml:"data,omitempty>mails>mail"`
}

type MetaData struct {
	Code    int    `xml:"code"`
	Message string `xml:"message"`
}

type MailData struct {
	EMailAddress string `xml:"emailAddress"`
	Subject      string `xml:"subject"`
	Message      string `xml:"message"`
}

// returns error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithXML(w, code, Response{Meta: MetaData{Code: code, Message: msg}})
}

// writes xml response
func RespondWithXML(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := xml.Marshal(payload)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
