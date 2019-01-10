package webserver

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCreateSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	CreateSessionCookie(rr, "TestID")

	assert.Equal(t, 1, len(rr.Result().Cookies()))
	assert.Equal(t, "TestID", rr.Result().Cookies()[0].Value)
}

func TestDestroySession(t *testing.T) {
	rr := httptest.NewRecorder()
	DestroySession(rr)

	assert.Equal(t, 1, len(rr.Result().Cookies()))
	assert.Equal(t, "session-id", rr.Result().Cookies()[0].Name)
	assert.Equal(t, "", rr.Result().Cookies()[0].Value)
}
