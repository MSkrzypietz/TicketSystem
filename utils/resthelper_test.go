package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithXMLFailure(t *testing.T) {
	payload := make(chan int)
	rr := httptest.NewRecorder()
	RespondWithXML(rr, http.StatusOK, payload)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/xml", rr.Header().Get("Content-Type"))
}

func TestRespondWithXMLSuccess(t *testing.T) {
	payload := Response{Meta: MetaData{Code: http.StatusOK, Message: "OK"}}
	rr := httptest.NewRecorder()
	RespondWithXML(rr, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/xml", rr.Header().Get("Content-Type"))
}
