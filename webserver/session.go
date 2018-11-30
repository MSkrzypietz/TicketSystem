package webserver

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"TicketSystem/utils"
	"net/http"
	"os"
)

func StartSession(w http.ResponseWriter, username string) {
	f, err := os.OpenFile("webserver/session_id.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cookieId := utils.CreateUUID(64)
	if _, err = f.WriteString(username + "," + cookieId); err != nil {
		panic(err)
	}
	CreateCookie(w, cookieId)
}

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

	return XML_IO.GetUserBySession(config.UsersPath, sessionID), err
}
