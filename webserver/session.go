package webserver

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

func StartSession(w http.ResponseWriter, username string) {
	f, err := os.OpenFile("webserver/session_id.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cookieId := CreateUUID(64)
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
	fmt.Fprintln(w, "Cookie set")
}

func DestroySession(r *http.Request) {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		panic(err)
	}
	cookie.Name = "Deleted"
	cookie.Value = "Unused"
	cookie.MaxAge = -1
}

func GetUserFromCookie(r *http.Request) (XML_IO.User, error) {
	sessionID := ""
	cookie, err := r.Cookie("session-id")
	if err == nil {
		sessionID = cookie.Value
	}

	return XML_IO.GetUserBySession(config.UsersPath, sessionID), err
}

func CreateUUID(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	byteSlice := make([]byte, length)
	for i := range byteSlice {
		byteSlice[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(byteSlice)
}
