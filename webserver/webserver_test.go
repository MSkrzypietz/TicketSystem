package webserver

import (
	"TicketSystem/XML_IO"
	"bytes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func createUser(userName string, password string) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return false, err
	}

	_, errUser := XML_IO.CreateUser(userName, string(hashedPassword))
	if errUser != nil {
		return false, err
	}

	return true, nil
}
