package webserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckUser(t *testing.T) {
	realUser := CheckUser("Test", "123456")
	unrealUser := CheckUser("Test", "12345")

	assert.True(t, realUser)
	assert.False(t, unrealUser)
}

func TestCreateUUID(t *testing.T) {
	uuid1 := CreateUUID(64)
	uuid2 := CreateUUID(64)
	assert.NotEqual(t, uuid1, uuid2)
	assert.Equal(t, 64, len(uuid1))
	assert.Equal(t, 64, len(uuid2))
}

func TestRealUser(t *testing.T) {
	assert.True(t, RealUser("Test"))
	assert.False(t, RealUser("FalseUser"))
}
