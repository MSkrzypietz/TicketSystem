package utils

import "strconv"

type Error int

const (
	ErrorUnknown Error = iota
	ErrorUnauthorized
	ErrorTemplateExecution
	ErrorInvalidTicketID
	ErrorDataFetching
	ErrorFormParsing
	ErrorInvalidInputs
	ErrorTicketCreation
	ErrorUserCreation
	ErrorUserLogin
	ErrorURLParsing
	ErrorDataStoring
)

// This is inspired by http://golang-basic.blogspot.com/2014/07/enumeration-example-golang.html
// This helps to prevent bugs by not having to manually maintaining a constant for the error count
var errors = []string{
	"We ran into an unknown error!",
	"You are not authorized for your actions!",
	"We had problems serving your requested site for you!",
	"The requested ticket does not exist!",
	"We had internal problems fetching the requested data!",
	"We had problems parsing your posted form!",
	"Your Inputs are invalid. Please check your inputs and try again!",
	"We had issues creating your ticket. Please try it again!",
	"We had issues signing you up to our services. Please try it again!",
	"We had issues signing you in. Please try it again!",
	"We had issues parsing your URL. Please try it again!",
	"We had issues storing your changes. Please try it again!",
}

// Returns the error message for a particular error
func (err Error) String() string {
	return errors[err]
}

// Returns the corresponding URL for a particular error
func (err Error) ErrorPageURL() string {
	return "/error/" + strconv.Itoa(int(err))
}

// Returns the total count of registered errors
func ErrorCount() int {
	return len(errors)
}
