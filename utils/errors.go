package utils

import "strconv"

type Error int

const (
	ErrorUnauthorized Error = iota + 1
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

func (err Error) String() string {
	switch err {
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
	return "Unkown Error"
}

func (err Error) ErrorPageUrl() string {
	return "/error/" + strconv.Itoa(int(err))
}
