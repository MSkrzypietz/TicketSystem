package main

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"TicketSystem/webserver"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"testing"
)

func setup() {
	config.DataPath = "datatest"
	config.TemplatePath = path.Join("..", "..", "templates")
	config.ServerKeyPath = path.Join("..", "..", "etc", "server.key")
	config.ServerCertPath = path.Join("..", "..", "etc", "server.crt")
	config.Port = getAvailablePort()
	webserver.Setup()
}

func getAvailablePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1
	}

	availablePort := listener.Addr().(*net.TCPAddr).Port
	if listener.Close() != nil {
		return -1
	}

	return availablePort
}

func teardown() {
	err := os.RemoveAll(config.DataPath)
	if err != nil {
		log.Println(err)
	}
}

func TestPullEmailsInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	shutdown := make(chan bool)
	done := make(chan bool)
	go webserver.StartServer(done, shutdown)

	_ = utils.SendMail("test@gmail.de", "Test Subject 1", "Test Message 1")
	_ = utils.SendMail("test@gmail.de", "Test Subject 2", "Test Message 2")
	_ = utils.SendMail("test@gmail.de", "Test Subject 3", "Test Message 3")

	emails, err := pullEmails("https://host:443")
	assert.NotNil(t, err)
	assert.Nil(t, emails)

	done <- true
	<-shutdown
}

func TestPullEmailsSuccess(t *testing.T) {
	setup()
	defer teardown()

	shutdown := make(chan bool)
	done := make(chan bool)
	go webserver.StartServer(done, shutdown)

	_ = utils.SendMail("test@gmail.de", "Test Subject 1", "Test Message 1")
	_ = utils.SendMail("test@gmail.de", "Test Subject 2", "Test Message 2")
	_ = utils.SendMail("test@gmail.de", "Test Subject 3", "Test Message 3")

	emails, err := pullEmails("https://localhost:" + strconv.Itoa(config.Port))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(emails))

	done <- true
	<-shutdown
}
