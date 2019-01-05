package webserver

import (
	"TicketSystem/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func setup() {
	config.DataPath = path.Join("..", "datatest")
	config.TemplatePath = path.Join("..", "templates")
	Setup()
}

func TestServeNewTicket(t *testing.T) {
	setup()
	defer os.RemoveAll(config.DataPath)

	req, err := http.NewRequest("POST", "/tickets/new", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeNewTicket)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeIndex(t *testing.T) {
	setup()
	defer os.RemoveAll(config.DataPath)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeIndex)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeSignOut(t *testing.T) {
	setup()
	defer os.RemoveAll(config.DataPath)

	req, err := http.NewRequest("POST", "/signOut", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignOut)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)

	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "session-id" || cookie.Name == "requested-url-while-not-authenticated" {
			assert.Equal(t, -1, cookie.MaxAge)
		}
	}
}

/*
func TestServeUserRegistration(t *testing.T) {
	if ok, err := createUser("Test123", "123"); !ok {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"newUsername": "Test123", "newPassword1"": 123, "newPassword2"": 123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		assert.Equal(t, status, http.StatusAccepted)
	}

	// TODO: Check if user is registered / available in users.xml with xmlIO.CheckUser

	// TODO: Remove User or delete users.xml file??
}
*/
