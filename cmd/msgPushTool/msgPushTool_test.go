package main

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"TicketSystem/config"
	"TicketSystem/webserver"
	"log"
	"os"
	"path"
)

func setup() {
	config.DataPath = "datatest"
	config.TemplatePath = path.Join("..", "..", "templates")
	config.ServerKeyPath = path.Join("..", "..", "etc", "server.key")
	config.ServerCertPath = path.Join("..", "..", "etc", "server.crt")
	webserver.Setup()
}

func teardown() {
	err := os.RemoveAll(config.DataPath)
	if err != nil {
		log.Println(err)
	}
}

/* Did not work with push ci...
func TestPushEmail(t *testing.T) {
	setup()
	defer teardown()

	shutdown := make(chan bool)
	done := make(chan bool)
	go webserver.StartServer(done, shutdown)

	email, err := pushEmail("https://localhost:"+strconv.Itoa(config.Port), "test@gmail.com", "Test Subject", "Test Message")
	assert.Nil(t, err)
	assert.NotNil(t, email)

	done <- true
	<-shutdown
}
*/
