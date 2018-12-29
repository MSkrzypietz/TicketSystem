package config

import "path"

var (
	DataPath       = "data"
	UsersFilePath  = path.Join(DataPath, "users", "users.xml")
	TicketsPath    = path.Join(DataPath, "tickets")
	ServerCertPath = path.Join("etc", "server.crt")
	ServerKeyPath  = path.Join("etc", "server.key")
	TemplatePath   = "templates"
	Port           = 4443
)
