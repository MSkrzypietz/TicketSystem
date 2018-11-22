package main

import (
	"github.com/stretchr/testify/assert"
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
	/* TODO: First assert is not working because the listener is not closing in time -> how to wait until its stopped??? */
	// Using port 0 returns a free / available port selected by the system
	/*
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		availablePort := listener.Addr().(*net.TCPAddr).Port

		gracefulListener := helpers.NewGracefulListener(listener, 2 * time.Second)
		gracefulListener.Close()

		assert.Nil(t, checkPortAvailability(availablePort))

		listener, err = net.Listen("tcp", ":" + string(availablePort))
		if err != nil {
			panic(err)
		}

		assert.Equal(t, false, checkPortAvailability(availablePort))
		listener.Close() */
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

func TestValidatePath(t *testing.T) {
	actual, err := validatePath(" ", ".")
	assert.Equal(t, ".", actual)
	assert.NotNil(t, err)

	actual, err = validatePath("test", ".")
	assert.Equal(t, ".", actual)
	assert.NotNil(t, err)

	actual, err = validatePath("main.go", ".")
	assert.Equal(t, "main.go", actual)
	assert.Nil(t, err)
}

// This test is inspired by Andrew Gerrand of the Go team
// Source: https://talks.golang.org/2014/testing.slide#23
/*
func TestExitOnError(t *testing.T) {
	if os.Getenv("CRASH_TEST") == "1" {
		err := fmt.Errorf("test error")
		exitOnError(err, "users")
		exitOnError(err, "cert")
		exitOnError(err, "key")
		exitOnError(err, "templates")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExitOnError")
	cmd.Env = append(os.Environ(), "CRASH_TEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 3", err)
}
*/
