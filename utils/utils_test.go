package utils

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckMailFormal(t *testing.T) {
	tests := []struct {
		mail     string
		expected bool
	}{
		{"ç$€§/az@gmail.com", false},
		{"abcd@gmail_yahoo.com", false},
		{"abcd@gmailyahoo.de", true},
		{"abcd@gmail.yahoo", true},
		{"", false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckMailFormal(d.mail))
	}
}

func TestCheckEmptyXSSString(t *testing.T) {
	longText := "asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi" +
		"asderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturiasderpotlamiutncturi too long"

	tests := []struct {
		text     string
		expected bool
	}{
		{"", false},
		{"<script>alert('XSS')</>", false},
		{"This is a Message", true},
		{longText, false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckEmptyXSSString(d.text))
	}
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
	tests := []struct {
		a        string
		expected bool
	}{
		{"test", false},
		{"test1", false},
		{"Test1", false},
		{"12Test!", true},
		{"1Te!", false},
		{"12 Test!", false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckPasswordFormal(d.a))
	}
}

func TestCheckEqualStrings(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		{"", "", false},
		{"hello", "test", false},
		{"test", "test", true},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckEqualStrings(d.a, d.b))
	}
}

func TestGetUserFromCookieSessionNotSet(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)

	_, err := GetUserFromCookie(req)
	assert.NotNil(t, err)
}

func TestGetUserFromCookieSuccess(t *testing.T) {
	setup()
	defer teardown()

	user1, err := CreateUser("Test123", "TestPassword")
	assert.Nil(t, err)
	user1.SessionID = "TestSession"
	assert.Nil(t, LoginUser(user1.Username, "TestPassword", user1.SessionID))

	req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    user1.SessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	user2, err := GetUserFromCookie(req)
	assert.Nil(t, err)
	assert.Equal(t, user1.Username, user2.Username)
}

func TestRemoveCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	RemoveCookie(rr, "session-id")

	assert.Equal(t, 1, len(rr.Result().Cookies()))
	assert.Equal(t, "session-id", rr.Result().Cookies()[0].Name)
	assert.Equal(t, "", rr.Result().Cookies()[0].Value)
}
