package main

import (
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/config"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/server"
	"log"
)

func main() {
	config.GetFlags()

	// Connect to the database
	database.Connect()

	// Start a bot specified from the configuration
	proxy := bot.GetProxy(bot.NewBot(config.GetServerConfiguration().Bot))

	go proxy.Start()

	// Start monitoring server
	err := server.Start(proxy.CommandsChannel, proxy.ReplyChannel)
	if err != nil {
		log.Println(err)
		return
	}
}
