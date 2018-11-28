package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegExMail(t *testing.T) {
	mail1 := "ç$€§/az@gmail.com"
	mail2 := "abcd@gmail_yahoo.com"
	mail3 := "abcd@gmailyahoo.de"
	mail4 := "abcd@gmail.yahoo"
	mail5 := ""

	assert.False(t, RegExMail(mail1))
	assert.False(t, RegExMail(mail2))
	assert.True(t, RegExMail(mail3))
	assert.True(t, RegExMail(mail4))
	assert.False(t, RegExMail(mail5))
}

func TestCheckEmptyXssString(t *testing.T) {
	text1 := ""
	text2 := "<script>alert('XSS')</>"
	text3 := "This is a Message"

	assert.False(t, CheckEmptyXssString(text1))
	assert.False(t, CheckEmptyXssString(text2))
	assert.True(t, CheckEmptyXssString(text3))
}
