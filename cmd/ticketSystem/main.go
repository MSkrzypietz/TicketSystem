package main

import (
	"TicketSystem/config"
	"TicketSystem/webserver"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strings"
)

func main() {
	handleFlags()

	webserver.StartServer()
}

func handleFlags() {
	usersFilePath := flag.String("users", "data/users/users.xml", "Path to users.xml file")
	serverCertPath := flag.String("cert", "etc/server.crt", "Path to server certificate")
	serverKeyPath := flag.String("key", "etc/server.key", "Path to server key")
	templatePath := flag.String("templates", "templates", "Path to templates folder")
	port := flag.Int("port", 443, "Port on which the server should run")
	flag.Parse()

	config.Port = *port
	if ok, err := validatePort(config.Port); !ok {
		log.Println(err)
		config.Port = 443

		/*
			if err := checkPortAvailability(config.Port); err != nil {
				log.Println("the default port 443 is also in use and hence the program will terminate")
				os.Exit(0)
			}
		*/
	}

	var err error

	config.UsersPath, err = validatePath(*usersFilePath, "data/users/users.xml")
	exitOnError(err, "users")

	config.ServerCertPath, err = validatePath(*serverCertPath, "https/server.crt")
	exitOnError(err, "cert")

	config.ServerKeyPath, err = validatePath(*serverKeyPath, "https/server.key")
	exitOnError(err, "key")

	config.TemplatePath, err = validatePath(*templatePath, "templates")
	exitOnError(err, "templates")
}

func exitOnError(err error, causer string) {
	if err != nil {
		log.Fatal(fmt.Errorf("Value of flag "+causer+" is invalid: %v", err))
		os.Exit(1)
	}
}

func validatePort(port int) (bool, error) {
	if !checkPortBoundaries(port) {
		return false, fmt.Errorf("the specified port %d is not in the port range of 0 to 65535. Using the default port 443 instead", port)
	}

	/* TODO: not working properly -> see test
	if err := checkPortAvailability(port); err != nil {
		return false, fmt.Errorf("the specified port %q is already in use: %s", port, err)
	}
	*/

	return true, nil
}

func checkPortBoundaries(port int) bool {
	return port >= 0 && port <= math.MaxUint16
}

func checkPortAvailability(port int) error {
	l, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		return err
	}
	defer l.Close()
	return nil
}

func validatePath(path string, def string) (string, error) {
	if !isPathSpecified(path) {
		return def, fmt.Errorf("no path specified")
	}

	if ok, err := existsPath(path); !ok {
		return def, fmt.Errorf("the specified path does not exist: %v", err)
	}

	return path, nil
}

func isPathSpecified(path string) bool {
	return strings.TrimSpace(path) != ""
}

func existsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err), err
}
