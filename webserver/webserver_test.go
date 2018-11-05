package webserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexPage(t *testing.T) {

}

func TestExecuteTemplate(t *testing.T) {
	//response := ExecuteTemplate(http.ResponseWriter(), *http.Request{}, "../templates/login.html", nil)
}

func TestReadTxtFile(t *testing.T) {
	result, err := ReadTxtFile("users.txt")
	print(result)
	print(result[1])
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
