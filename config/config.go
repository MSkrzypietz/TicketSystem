package config

import "path"

var (
	DataPath       = "data"
	ServerCertPath = path.Join("etc", "server.crt")
	ServerKeyPath  = path.Join("etc", "server.key")
	TemplatePath   = "templates"
	Port           = 4443
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
