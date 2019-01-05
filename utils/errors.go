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

const ErrorCount = 12

func (err Error) String() string {
	switch err {
	case ErrorUnknown:
		return "We ran into an unknown error!"
	case ErrorUnauthorized:
		return "You are not authorized for your actions!"
	case ErrorTemplateExecution:
		return "We had problems serving your requested site for you!"
	case ErrorInvalidTicketID:
		return "The requested ticket does not exist!"
	case ErrorDataFetching:
		return "We had internal problems fetching the requested data!"
	case ErrorFormParsing:
		return "We had problems parsing your posted form!"
	case ErrorInvalidInputs:
		return "Your Inputs are invalid. Please check your inputs and try again!"
	case ErrorTicketCreation:
		return "We had issues creating your ticket. Please try it again!"
	case ErrorUserCreation:
		return "We had issues signing you up to our services. Please try it again!"
	case ErrorUserLogin:
		return "We had issues signing you in. Please try it again!"
	case ErrorURLParsing:
		return "We had issues parsing your URL. Please try it again!"
	case ErrorDataStoring:
		return "We had issues storing your changes. Please try it again!"
	}
	return ErrorUnknown.String()
}

func (err Error) ErrorPageUrl() string {
	return "/error/" + strconv.Itoa(int(err))
}
