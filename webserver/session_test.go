package webserver

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCreateSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	createSessionCookie(rr, "TestID")

	assert.Equal(t, 1, len(rr.Result().Cookies()))
	assert.Equal(t, "TestID", rr.Result().Cookies()[0].Value)
}

func TestDestroySession(t *testing.T) {
	rr := httptest.NewRecorder()
	destroySession(rr)

	assert.Equal(t, 1, len(rr.Result().Cookies()))
	assert.Equal(t, "session-id", rr.Result().Cookies()[0].Name)
	assert.Equal(t, "", rr.Result().Cookies()[0].Value)
}
