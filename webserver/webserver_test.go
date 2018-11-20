package webserver

import (
	"TicketSystem/config"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage(t *testing.T) {

}

func TestExecuteTemplate(t *testing.T) {
	//response := ExecuteTemplate(http.ResponseWriter(), *http.Request{}, "../templates/login.html", nil)
}

func TestServeUserRegistration(t *testing.T) {
	config.UsersPath = "../data/users/users.xml"

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"newUsername": "Test123", "newPassword1"": 123, "newPassword2"": 123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned the wrong status code... got %v and expected %v", status, http.StatusAccepted)
	}

	// TODO: Check if user is registered / available in users.xml with xmlIO.CheckUser

	// TODO: Remove User or delete users.xml file??
}
