package config

type Configuration struct {
	Server   Server
	Telegram Telegram
	Mongo    Database
}

type Server struct {
	Host         string `fig:"host" validate:"required"`
	Port         int    `fig:"port" default:"7777"`
	AuthPassword string `fig:"authPassword" validate:"required"`
	PluginDir    string `fig:"host" validate:"required"`
	Plugins      []string
}

type Telegram struct {
	BotToken string `fig:"host" validate:"required"`
}

type Database struct {
	Scheme   string
	Host     string `fig:"host" default:"localhost"`
	Username string `fig:"host" validate:"required"`
	Password string `fig:"host" validate:"required"`
	Port     int    `fig:"host" default:"2717"`
	Database string `fig:"host" default:"sn-bot"`
}
