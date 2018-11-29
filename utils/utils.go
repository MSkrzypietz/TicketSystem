package utils

import (
	"math/rand"
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
