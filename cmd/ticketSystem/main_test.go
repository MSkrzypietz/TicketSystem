package main

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
)

func TestCheckPortBoundaries(t *testing.T) {
	assert.Equal(t, true, checkPortBoundaries(0))
	assert.Equal(t, true, checkPortBoundaries(1337))
	assert.Equal(t, true, checkPortBoundaries(65535))

	assert.Equal(t, false, checkPortBoundaries(-1))
	assert.Equal(t, false, checkPortBoundaries(65536))
}

func TestCheckPortAvailability(t *testing.T) {
	// Using port 0 returns a free / available port selected by the system
	listener, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)
	availablePort := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	assert.Nil(t, err)

	actual, err := checkPortAvailability(availablePort)
	assert.True(t, actual)
	assert.Nil(t, err)

	listener, err = net.Listen("tcp", ":"+strconv.Itoa(availablePort))
	assert.Nil(t, err)
	defer listener.Close()

	actual, err = checkPortAvailability(availablePort)
	assert.False(t, actual)
	assert.NotNil(t, err)
}

func TestIsPathSpecified(t *testing.T) {
	assert.False(t, isPathSpecified(""))
	assert.False(t, isPathSpecified(" "))

	assert.True(t, isPathSpecified("."))
	assert.True(t, isPathSpecified("test"))
}

func TestExistsPath(t *testing.T) {
	actual, err := existsPath("test")
	assert.False(t, actual)
	assert.NotNil(t, err)

	actual, err = existsPath("main.go")
	assert.True(t, actual)
	assert.Nil(t, err)

	actual, err = existsPath("../ticketSystem/main.go")
	assert.True(t, actual)
	assert.Nil(t, err)
}

func TestCheckPath(t *testing.T) {
	var testVar *string
	ok, err := checkPath(path{" ", testVar, "test"})
	assert.False(t, ok)
	assert.NotNil(t, err)

	ok, err = checkPath(path{"test", testVar, "test"})
	assert.False(t, ok)
	assert.NotNil(t, err)

	ok, err = checkPath(path{"main.go", testVar, "test"})
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestHandlePort(t *testing.T) {

}
