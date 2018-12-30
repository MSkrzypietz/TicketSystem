package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"net/http"
)

func CreateCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session-id",
		Value:  id,
		MaxAge: 60 * 60,
	})
	//log.Println(w, "Cookie set")
}

func DestroySession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func GetUserFromCookie(r *http.Request) (XML_IO.User, error) {
	sessionID := ""
	cookie, err := r.Cookie("session-id")
	if err == nil {
		sessionID = cookie.Value
	}

	return XML_IO.GetUserBySession(config.UsersFilePath(), sessionID), err
}
