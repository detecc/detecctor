package config

type (
	Configuration struct {
		Server Server
		Bot    BotConfiguration
		Mongo  Database
	}
	Server struct {
		Host         string `fig:"host"`
		Port         int    `fig:"port" default:"7777"`
		AuthPassword string `fig:"authPassword" validate:"required"`
		PluginDir    string `fig:"pluginDir" validate:"required"`
		Plugins      []string
	}
	BotConfiguration struct {
		Token string `fig:"token" validate:"required"`
		Type  string `fig:"type" validate:"required"`
		ID    string `fig:"id" validate:"required"`
	}

	Database struct {
		Scheme   string
		Host     string `fig:"host" default:"localhost"`
		Username string `fig:"username" validate:"required"`
		Password string `fig:"password" validate:"required"`
		Port     int    `fig:"port" default:"2717"`
		Database string `fig:"database" default:"detecctor"`
	}
)
