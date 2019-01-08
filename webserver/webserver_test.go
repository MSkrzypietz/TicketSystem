package webserver

import (
	"TicketSystem/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticateWithoutCookie(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/tickets/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate(ServeTickets))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusFound, rr.Code)

	// RemoveCookie counts as another cookie even tho it will get destroyed instantly; hence -> expected: 2
	assert.Equal(t, 2, len(rr.Result().Cookies()))

	location, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/signIn", location.Path)
}

func TestAuthenticateWithCookieButInvalidSession(t *testing.T) {
	setup()
	defer teardown()

	username := "Test123"
	password := "Aa!123456"
	uuid := utils.CreateUUID(64)

	createUser(username, password)
	err := utils.LoginUser(username, password, uuid)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/tickets/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    "wrongSessionID",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate(ServeTickets))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusFound, rr.Code)

	// RemoveCookie counts as another cookie even tho it will get destroyed instantly; hence -> expected: 2
	assert.Equal(t, 2, len(rr.Result().Cookies()))

	location, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/signIn", location.Path)
}
