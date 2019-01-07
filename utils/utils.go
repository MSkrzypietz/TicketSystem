package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

func RegExMail(email string) bool {
	emailRegExp := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
			"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
			"[a-zA-Z0-9])?)*$")
	return emailRegExp.MatchString(email)
}

func CheckEmptyXssString(text string) bool {
	if len(text) == 0 || len(text) > 400 {
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

// Checks if the inputs contains only ASCII letters and digits, with hyphens, underscores and spaces
// as internal separators.
// Source: https://stackoverflow.com/questions/1221985/how-to-validate-a-user-name-with-regex
func CheckUsernameFormal(name string) bool {
	if len(name) <= 4 || len(name) >= 30 {
		return false
	}

	nameRegExp := regexp.MustCompile("^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$")
	return nameRegExp.MatchString(name)
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

func GetUserFromCookie(r *http.Request) (User, error) {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		return User{}, fmt.Errorf("session id is not set")
	}

	return GetUserBySession(cookie.Value)
}

func RemoveCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
