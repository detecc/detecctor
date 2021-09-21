package main

import (
	"log"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/config"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/server"
)

func main() {
	config.GetFlags()

	// Connect to the database
	database.Connect()

	// Start telegram bot
	telegram := bot.NewBot(config.GetServerConfiguration().Telegram.BotToken)
	go telegram.Start()

	// Start monitoring server
	err := server.Start(telegram.CommandsChannel, telegram.ReplyChannel)
	if err != nil {
		log.Println(err)
		return
	}
}
