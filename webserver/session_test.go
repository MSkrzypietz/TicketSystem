package webserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUUID(t *testing.T) {
	uuid1 := CreateUUID(64)
	uuid2 := CreateUUID(64)
	assert.NotEqual(t, uuid1, uuid2)
	assert.Equal(t, 64, len(uuid1))
	assert.Equal(t, 64, len(uuid2))
}


func TestGetUserFromCookie(t *testing.T) {

}