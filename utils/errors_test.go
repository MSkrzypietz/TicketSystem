package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		err      Error
		expected string
	}{
		{ErrorDataStoring, errors[ErrorDataStoring]},
		{ErrorUnauthorized, errors[ErrorUnauthorized]},
		{ErrorInvalidTicketID, errors[ErrorInvalidTicketID]},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, d.err.String())
	}
}

func TestErrorPageURL(t *testing.T) {
	tests := []struct {
		err      Error
		expected string
	}{
		{ErrorDataStoring, "/error/" + strconv.Itoa(int(ErrorDataStoring))},
		{ErrorUnauthorized, "/error/" + strconv.Itoa(int(ErrorUnauthorized))},
		{ErrorInvalidTicketID, "/error/" + strconv.Itoa(int(ErrorInvalidTicketID))},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, d.err.ErrorPageURL())
	}
}

func TestErrorCount(t *testing.T) {
	assert.Equal(t, len(errors), ErrorCount())
}
