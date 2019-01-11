package utils

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

// Checks if the email is valid
func CheckMailFormal(email string) bool {
	emailRegExp := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
			"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
			"[a-zA-Z0-9])?)*$")
	return emailRegExp.MatchString(email)
}

// Checks if the string is empty, too long and doesn't contain XSS
func CheckEmptyXSSString(text string) bool {
	if len(text) == 0 || len(text) > 400 {
		return false
	}

	xss := "<>[]"
	return !strings.ContainsAny(text, xss)
}

// Creates an universally unique identifier
func CreateUUID(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	byteSlice := make([]byte, length)
	for i := range byteSlice {
		byteSlice[i] = letters[rand.Int63()%int64(len(letters))]
	}

	return string(byteSlice)
}

// Checks if the inputs contains only ASCII letters and digits, with hyphens, underscores and spaces
// as internal separators. Username needs at least 5 chars, it can't exceed 29 chars
// Source: https://stackoverflow.com/questions/1221985/how-to-validate-a-user-name-with-regex
func CheckUsernameFormal(name string) bool {
	if len(name) <= 4 || len(name) >= 30 {
		return false
	}

	nameRegExp := regexp.MustCompile("^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$")
	return nameRegExp.MatchString(name)
}

// Checks safety of a password by requiring numbers, letters, uppercase letters and symbols in the passwords.
// Spaces are forbidden
func CheckPasswordFormal(password string) bool {
	if len(password) <= 4 {
		return false
	}

	numbers := "0123456789"
	letters := "abcdefghijklmnopqrstuvwxyz"
	symbols := "!#+-*.,:;/()§$%&?"
	space := " "

	if strings.ContainsAny(password, numbers) && strings.ContainsAny(password, letters) &&
		strings.ContainsAny(password, strings.ToUpper(letters)) && strings.ContainsAny(password, symbols) &&
		!strings.ContainsAny(password, space) {
		return true
	}

	return false
}

// Checks if strings are not empty and equal
func CheckEqualStrings(string1, string2 string) bool {
	return string1 != "" && string2 != "" && string1 == string2
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
