package webserver

import (
	"TicketSystem/utils"
	"net/http"
)

func CreateSessionCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    id,
		HttpOnly: true,
		MaxAge:   60 * 60,
	})
}

func DestroySession(w http.ResponseWriter) {
	utils.RemoveCookie(w, "session-id")
}
