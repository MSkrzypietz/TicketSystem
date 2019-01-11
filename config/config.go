package config

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"path"
	"strconv"
)

var (
	DataPath       = "data"
	ServerCertPath = path.Join("etc", "server.crt")
	ServerKeyPath  = path.Join("etc", "server.key")
	TemplatePath   = "templates"
	Port           = 4443
	DebugMode      = true
)

func UsersPath() string {
	return path.Join(DataPath, "users")
}

func UsersFilePath() string {
	return path.Join(UsersPath(), "users.xml")
}

func TicketsPath() string {
	return path.Join(DataPath, "tickets")
}

func TicketXMLPath(id int) string {
	return path.Join(TicketsPath(), "ticket"+strconv.Itoa(id)+".xml")
}

func DefinitionsFilePath() string {
	return path.Join(DataPath, "definitions.xml")
}

func MailFilePath() string {
	return path.Join(DataPath, "mails.xml")
}
