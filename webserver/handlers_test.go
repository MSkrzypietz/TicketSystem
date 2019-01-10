package webserver

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TODO: In jeder Testdatei anstelle assert.NotNil(t, err) -> assert.NoError(t, err)

func setup() {
	config.DataPath = "datatest"
	config.TemplatePath = path.Join("..", "templates")
	Setup()
}

func teardown() {
	err := os.RemoveAll(config.DataPath)
	if err != nil {
		log.Println(err)
	}
}

func TestServeNewTicket(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeNewTicket)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeIndex(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeIndex)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeErrorPage(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/error/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "/error/1", req.URL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error/1000", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err = rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error/1test", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err = rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)
}

func TestServeSignOut(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/signOut", nil)
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

func TestServeUserRegistrationShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "")
	form.Add("password2", "123")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestServeUserRegistrationPasswordsDontMatch(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "123456")
	form.Add("password2", "12345")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeUserRegistrationInvalidUsername(t *testing.T) {
	setup()
	defer teardown()

	config.DebugMode = false

	form := url.Values{}
	form.Add("username", "Tes")
	form.Add("password1", "12345")
	form.Add("password2", "12345")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeUserRegistrationInvalidPassword(t *testing.T) {
	setup()
	defer teardown()

	config.DebugMode = false

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "123")
	form.Add("password2", "123")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeUserRegistrationSuccess(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "Aa!123456")
	form.Add("password2", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeAuthenticationShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func createUser(username, password string) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("password1", password)
	form.Add("password2", password)

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)
}

func TestServeAuthenticationInvalidCredentials(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Tes23")
	form.Add("password", "A3456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUserLogin.ErrorPageURL(), resultURL.Path)
}

func TestServeAuthenticationSuccessWithHomeRedirect(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test123", "Aa!123456")

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeAuthenticationSuccessWithRequestedURLRedirect(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test123", "Aa!123456")
	requestedURL := "/ticket/new"

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form
	req.AddCookie(&http.Cookie{
		Name:     "requested-url-while-not-authenticated",
		Value:    requestedURL,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, requestedURL, resultURL.Path)
}

func TestServeTicketCreationInvalidInputs(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("email", "mustermann@gmail.com")
	form.Add("subject", "PC Issue")

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketCreation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketCreationSuccess(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("email", "mustermann@gmail.com")
	form.Add("subject", "PC Issue")
	form.Add("message", "I have issues with my pc...")

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketCreation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeAddCommentUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "Test Comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentInvalidInput(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func loginUser(rr *httptest.ResponseRecorder, username, password, uuid string) error {
	err := utils.LoginUser(username, password, uuid)
	createSessionCookie(rr, uuid)
	return err
}

func TestServeAddCommentInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "My comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "My comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidTicketID.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentSuccessComment(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("comment", "My comment")
	form.Add("sendoption", "comments")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeAddCommentSuccessEmail(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("comment", "My comment")
	form.Add("sendoption", "customer")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func createDummyTicket() (utils.Ticket, error) {
	return utils.CreateTicket("test@gmail.com", "Subject Dummy", "Message dummy")
}

func TestServeTicketAssignmentUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidEditor(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorDataStoring.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentSuccessWithoutRedirect(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeTicketAssignmentSuccessWithRedirect(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test123456")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	createUser("Test123456", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/tickets/", resultURL.Path)
}

func TestServeTicketReleaseUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorDataFetching.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidUser(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test1234567", "Aa!123456")
	testTicket, err := createDummyTicket()
	assert.Nil(t, utils.ChangeEditor(testTicket.Id, "Test1234567"))
	assert.Nil(t, err)

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseSuccess(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, utils.ChangeEditor(testTicket.Id, "Test123"))
	assert.Nil(t, err)

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeTicketsUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/tickets/", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTickets)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketsShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/ticket/", nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeTicketsInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/ticket/1337", nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidTicketID.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketsSuccess(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test124563", "Aa!123456")

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/ticket/"+strconv.Itoa(testTicket.Id), nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeCloseTicketUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeCloseTicketInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeCloseTicketSuccess(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/ticket/"+strconv.Itoa(testTicket.Id))

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/tickets/", resultURL.Path)
}

func TestPreventEMailPingPong(t *testing.T) {
	tests := []struct {
		mail     utils.Mail
		expected bool
	}{
		{utils.Mail{ReadAttemptCounter: 0, FirstReadAttemptDate: time.Now().Add(-30 * time.Second)}, false},
		{utils.Mail{ReadAttemptCounter: 3, FirstReadAttemptDate: time.Now().Add(-5 * time.Second)}, false},
		{utils.Mail{ReadAttemptCounter: 4, FirstReadAttemptDate: time.Now().Add(-3 * time.Minute)}, true},
		{utils.Mail{ReadAttemptCounter: 4, FirstReadAttemptDate: time.Now().Add(-6 * time.Minute)}, false},
		{utils.Mail{ReadAttemptCounter: 6, FirstReadAttemptDate: time.Now().Add(-10 * time.Minute)}, true},
		{utils.Mail{ReadAttemptCounter: 6, FirstReadAttemptDate: time.Now().Add(-15 * time.Minute)}, false},
		{utils.Mail{ReadAttemptCounter: 9, FirstReadAttemptDate: time.Now().Add(-25 * time.Minute)}, true},
		{utils.Mail{ReadAttemptCounter: 9, FirstReadAttemptDate: time.Now().Add(-29 * time.Minute)}, false},
		{utils.Mail{ReadAttemptCounter: 15, FirstReadAttemptDate: time.Now().Add(-60 * time.Minute)}, true},
		{utils.Mail{ReadAttemptCounter: 15, FirstReadAttemptDate: time.Now().Add(-80 * time.Minute)}, false},
		{utils.Mail{ReadAttemptCounter: 20, FirstReadAttemptDate: time.Now().Add(-120 * time.Minute)}, true},
		{utils.Mail{ReadAttemptCounter: 20, FirstReadAttemptDate: time.Now().Add(-160 * time.Minute)}, false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, preventEMailPingPong(d.mail))
	}
}

func TestGetMailsFileReadError(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"

	req := httptest.NewRequest(http.MethodGet, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMails)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	config.DataPath = "datatest"
}

func TestGetMailsSuccess(t *testing.T) {
	setup()
	defer teardown()

	emailAddress := "test@gmail.com"
	subject := "Test Subject"
	message := "Test Message"
	err := utils.SendMail(emailAddress, subject, message)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodGet, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMails)

	for i := 0; i < 4; i++ {
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var mails utils.Response
		err = xml.NewDecoder(rr.Body).Decode(&mails)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(mails.Data))
	}

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var mails utils.Response
	err = xml.NewDecoder(rr.Body).Decode(&mails)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(mails.Data)) // EMail ping pong prevention strikes
}

func TestPostMailsInvalidPayload(t *testing.T) {
	setup()
	defer teardown()

	payload := `<Invalid Payload`

	req := httptest.NewRequest(http.MethodPost, "/mails", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postMails)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostMailsStorageFailure(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"

	mailReq := utils.Request{Mail: utils.MailData{EMailAddress: "Test@gmail.com", Subject: "Test Subject", Message: "Test Message"}}
	payload, err := xml.Marshal(mailReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mails", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postMails)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPostMailsSuccess(t *testing.T) {
	setup()
	defer teardown()

	mailReq := utils.Request{Mail: utils.MailData{EMailAddress: "Test@gmail.com", Subject: "Test Subject", Message: "Test Message"}}
	payload, err := xml.Marshal(mailReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mails", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postMails)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	_, err = utils.ReadTicket(1)
	assert.Nil(t, err) // ticket exists
}

func TestServeMailsAPIInvalidHTTPMethod(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodDelete, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestServeMailsAPIGetSuccess(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/mails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeMailsAPIPostSuccess(t *testing.T) {
	setup()
	defer teardown()

	mailReq := utils.Request{Mail: utils.MailData{EMailAddress: "Test@gmail.com", Subject: "Test Subject", Message: "Test Message"}}
	payload, err := xml.Marshal(mailReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mails", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsAPI)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeMailsSentNotificationInvalidHTTPMethod(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/mails/notify", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsSentNotification)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestServeMailsSentNotificationInvalidPayload(t *testing.T) {
	setup()
	defer teardown()

	payload := `<Invalid Payload`

	req := httptest.NewRequest(http.MethodPost, "/mails/notify", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsSentNotification)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestServeMailsSentNotificationInternalServerError(t *testing.T) {
	setup()
	defer teardown()

	notifyReq := utils.Request{MailIDs: []int{1}}
	payload, err := xml.Marshal(notifyReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mails/notify", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsSentNotification)

	config.DataPath = "wrongPath"
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	config.DataPath = "datatest"
}

func TestServeMailsSentNotificationSuccess(t *testing.T) {
	setup()
	defer teardown()

	err := utils.SendMail("Test@gmail.com", "Test Subject", "Test Message")
	assert.Nil(t, err)
	maillist, err := utils.ReadMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(maillist.Maillist))

	notifyReq := utils.Request{MailIDs: []int{1}}
	payload, err := xml.Marshal(notifyReq)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mails/notify", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeMailsSentNotification)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	maillist, err = utils.ReadMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(maillist.Maillist))
}
