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
	text4 := "asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi too long"

	assert.False(t, CheckEmptyXssString(text1))
	assert.False(t, CheckEmptyXssString(text2))
	assert.True(t, CheckEmptyXssString(text3))
	assert.False(t, CheckEmptyXssString(text4))
}

func TestCreateUUID(t *testing.T) {
	uuid1 := CreateUUID(64)
	uuid2 := CreateUUID(64)
	assert.NotEqual(t, uuid1, uuid2)
	assert.Equal(t, 64, len(uuid1))
	assert.Equal(t, 64, len(uuid2))
}

func TestCheckUsernameFormal(t *testing.T) {
	tests := []struct {
		a        string
		expected bool
	}{
		{"user", false},
		{"InvalidTestWithOverThirtyCharacters", false},
		{"_username", false},
		{"username_", false},
		{"us  ername", false},
		{"_us--ername", false},
		{"username", true},
		{"1username", true},
		{"username1", true},
		{"use1234rname", true},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckUsernameFormal(d.a))
	}
}

func TestCheckPasswdFormal(t *testing.T) {
	passwd1 := "test"     // not valid
	passwd2 := "test1"    // not valid
	passwd3 := "Test1"    // not valid
	passwd4 := "12Test!"  // valid
	passwd5 := "1Te!"     // not valid
	passwd6 := "12 Test!" // not valid

	assert.False(t, CheckPasswdFormal(passwd1))
	assert.False(t, CheckPasswdFormal(passwd2))
	assert.False(t, CheckPasswdFormal(passwd3))
	assert.True(t, CheckPasswdFormal(passwd4))
	assert.False(t, CheckPasswdFormal(passwd5))
	assert.False(t, CheckPasswdFormal(passwd6))
}

func TestCheckEqualStrings(t *testing.T) {
	str1 := ""
	str2 := "hello"
	str3 := "test"

	assert.False(t, CheckEqualStrings(str1, str1))
	assert.False(t, CheckEqualStrings(str2, str3))
	assert.True(t, CheckEqualStrings(str3, str3))
}
