package main

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"TicketSystem/config"
	"TicketSystem/webserver"
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	handleFlags()

	shutdown := make(chan bool)
	done := make(chan bool)
	go webserver.StartServer(done, shutdown)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Type quit to shut the ticket system down.")

		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		if strings.TrimSpace(strings.ToLower(input)) == "quit" {
			done <- true
			break
		}
	}

	// Waiting for server to gracefully shut down before exiting the program
	<-shutdown
}

func handleFlags() {
	flag.String("data", config.DataPath, "Path to data folder")
	serverCertPath := flag.String("cert", config.ServerCertPath, "Path to server certificate")
	serverKeyPath := flag.String("key", config.ServerKeyPath, "Path to server key")
	templatePath := flag.String("templates", config.TemplatePath, "Path to templates folder")
	port := flag.Int("port", config.Port, "Port on which the server should run")
	debugMode := flag.Bool("debug", config.DebugMode, "Decides the mode the server should run on")
	flag.Parse()

	if !checkPortBoundaries(*port) {
		log.Fatalf("Invalid port %d", *port)
	}
	handlePaths(*serverCertPath, *serverKeyPath, *templatePath)

	config.ServerCertPath = *serverCertPath
	config.ServerKeyPath = *serverKeyPath
	config.TemplatePath = *templatePath
	config.Port = *port
	config.DebugMode = *debugMode
}

func checkPortBoundaries(port int) bool {
	return port >= 0 && port <= math.MaxUint16
}

func handlePaths(paths ...string) {
	for _, path := range paths {
		if !existsPath(path) {
			log.Fatalf("Path %s does not exist!", path)
		}
	}
}

func existsPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
