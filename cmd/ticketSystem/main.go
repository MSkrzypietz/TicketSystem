package main

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"TicketSystem/config"
	"TicketSystem/webserver"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

type path struct {
	userInput string
	configVar *string
	varName   string
}

func main() {
	handleFlags()

	webserver.StartServer()
}

func handleFlags() {
	//TODO: Fix data path handling
	//dataPath := flag.String("data", config.DataPath, "Path to data folder")
	serverCertPath := flag.String("cert", config.ServerCertPath, "Path to server certificate")
	serverKeyPath := flag.String("key", config.ServerKeyPath, "Path to server key")
	templatePath := flag.String("templates", config.TemplatePath, "Path to templates folder")
	port := flag.Int("port", config.Port, "Port on which the server should run")
	debugMode := flag.Bool("debug", config.DebugMode, "Decides the mode the server should run on")
	flag.Parse()

	handlePort(*port)
	handlePathInputs(
		//path{*dataPath, &config.DataPath, "data"},
		path{*serverCertPath, &config.ServerCertPath, "cert"},
		path{*serverKeyPath, &config.ServerKeyPath, "key"},
		path{*templatePath, &config.TemplatePath, "templates"},
	)
	config.DebugMode = *debugMode
}

func handlePort(port int) {
	if ok, err := validatePort(port); !ok {
		log.Println(err)

		if ok, err := checkPortAvailability(config.Port); !ok {
			log.Fatalf("the default port %d is also in use. Hence the program will terminate: %v", config.Port, err)
		}

		log.Println(fmt.Errorf("the requested port is already in use. Using the default port %d instead", config.Port))
		return
	}

	config.Port = port
}

func validatePort(port int) (bool, error) {
	if !checkPortBoundaries(port) {
		return false, fmt.Errorf("the specified port %d is not in the port range of 0 to 65535. Using the default port 443 instead", port)
	}

	if ok, err := checkPortAvailability(port); !ok {
		return false, fmt.Errorf("the specified port %q is already in use: %s", port, err)
	}

	return true, nil
}

func checkPortBoundaries(port int) bool {
	return port >= 0 && port <= math.MaxUint16
}

func checkPortAvailability(port int) (bool, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false, err
	}
	err = l.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func handlePathInputs(paths ...path) {
	for _, path := range paths {
		if ok, err := checkPath(path); !ok {
			log.Fatal(fmt.Errorf("value of flag %s is invalid: %v", path.varName, err))
		}
		*path.configVar = path.userInput
	}
}

func checkPath(path path) (bool, error) {
	if !isPathSpecified(path.userInput) {
		return false, fmt.Errorf("no path specified")
	}

	if ok, err := existsPath(path.userInput); !ok {
		return false, fmt.Errorf("the specified path does not exist: %v", err)
	}

	return true, nil
}

func isPathSpecified(path string) bool {
	return strings.TrimSpace(path) != ""
}

func existsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err), err
}
