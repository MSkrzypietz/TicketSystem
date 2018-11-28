package utils

import (
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
	if text == "" {
		return false
	}
	xss := "<>[]"
	return !strings.ContainsAny(text, xss)
}
