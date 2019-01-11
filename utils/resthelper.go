package utils

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"encoding/xml"
	"log"
	"net/http"
)

// This file is inspired by https://itnext.io/building-restful-web-api-service-using-golang-chi-mysql-d85f427dee54

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

// Writes xml error response
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithXML(w, code, Response{Meta: MetaData{Code: code, Message: msg}})
}

// Writes xml response for any payload that can be converted into XML
func RespondWithXML(w http.ResponseWriter, code int, payload interface{}) {
	response, err := xml.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't marshal the response payload")
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
