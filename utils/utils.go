package utils

import (
	"TicketSystem/XML_IO"
	"TicketSystem/config"
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

func RegExMail(email string) bool {
	mailResEx := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
		"[a-zA-Z0-9])?)*$")
	return mailResEx.MatchString(email)
}

func CheckEmptyXssString(text string) bool {
	if text == "" || len(text) > 400 {
		return false
	}
	xss := "<>[]"
	return !strings.ContainsAny(text, xss)
}

func CreateUUID(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	byteSlice := make([]byte, length)
	for i := range byteSlice {
		byteSlice[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(byteSlice)
}

func CheckUsernameFormal(name string) bool {
	if name == "" || len(name) < 4 {
		return false
	}
	numbers := "0123456789"
	letters := "abcdefghijklmnopqrstuvwxyz"
	space := " "

	if strings.ContainsAny(name, numbers) {
		if strings.ContainsAny(name, letters) {
			return !strings.ContainsAny(name, space)
		} else {
			// No letters in username
			return false
		}
	} else {
		// No number in username
		return false
	}
}

func CheckPasswdFormal(passwd string) bool {
	if len(passwd) <= 4 {
		return false
	}
	numbers := "0123456789"
	letters := "abcdefghijklmnopqrstuvwxyz"
	symbols := "!#+-*.,:;/()§$%&?"
	space := " "

	if strings.ContainsAny(passwd, numbers) && strings.ContainsAny(passwd, letters) &&
		strings.ContainsAny(passwd, strings.ToUpper(letters)) && strings.ContainsAny(passwd, symbols) &&
		!strings.ContainsAny(passwd, space) {
		return true
	} else {
		// No number or capital or normal letter or symbol or whitespace in password
		return false
	}

}

func CheckEqualStrings(string1, string2 string) bool {
	// Check strings on empty strings and matching
	if string1 != "" && string2 != "" {
		if string1 == string2 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func GetUserFromCookie(r *http.Request) (XML_IO.User, error) {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		return XML_IO.User{}, errors.New("session id is not set")
	}

	return XML_IO.GetUserBySession(config.UsersFilePath(), cookie.Value)
}

func RemoveCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
