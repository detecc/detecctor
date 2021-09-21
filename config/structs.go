package config

type Configuration struct {
	Server   Server
	Telegram Telegram
	Mongo    Database
}

type Server struct {
	Host         string `fig:"host"`
	Port         int    `fig:"port" default:"7777"`
	AuthPassword string `fig:"authPassword" validate:"required"`
	PluginDir    string `fig:"pluginDir" validate:"required"`
	Plugins      []string
}

type Telegram struct {
	BotToken string `fig:"botToken" validate:"required"`
}

type Database struct {
	Scheme   string
	Host     string `fig:"host" default:"localhost"`
	Username string `fig:"username" validate:"required"`
	Password string `fig:"password" validate:"required"`
	Port     int    `fig:"port" default:"2717"`
	Database string `fig:"database" default:"detecctor"`
}
