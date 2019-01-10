package webserver

import (
	"TicketSystem/utils"
	"net/http"
)

func createSessionCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})
}

func destroySession(w http.ResponseWriter) {
	utils.RemoveCookie(w, "session-id")
}
