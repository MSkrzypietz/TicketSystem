package utils

import (
	"encoding/xml"
	"log"
	"net/http"
)

// this file is inspired by https://itnext.io/building-restful-web-api-service-using-golang-chi-mysql-d85f427dee54

type RestError struct {
	Code int
	Msg  string
}

type Response struct {
	Success bool      `xml:"success"`
	Data    []Mail    `xml:"data>mails>mail"`
	Err     RestError `xml:"error"`
}

type MailIDs struct {
	IDList []int `xml:"MailID"`
}

// returns error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithXML(w, code, Response{Success: false, Err: RestError{Code: code, Msg: msg}})
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
