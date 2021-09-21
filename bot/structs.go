package bot

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"github.com/detecc/detecctor/shared"
)

// Telegram is a wrapper for the Telegram bot API.
type Telegram struct {
	Token           string
	BotAPI          *telegram.BotAPI
	CommandsChannel chan Command
	ReplyChannel    chan shared.Reply
}

// Command consists of a Name and Args.
// Args contains the arguments extracted from the Message sent through Telegram.
// Example of a command: "/get_status serviceNode1 serviceNode2".
// The command name is "/get_status", the arguments are ["serviceNode1", "serviceNode2"].
type Command struct {
	Name   string
	Args   []string
	ChatId int64
}

// NewBot create a new telegram bot
func NewBot(token string) *Telegram {
	return &Telegram{
		Token:           token,
		CommandsChannel: make(chan Command),
		ReplyChannel:    make(chan shared.Reply),
	}
}

// Start listening to the bot updates and the updates from the TCP server
func (t *Telegram) Start() {
	bot, err := telegram.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}
	t.BotAPI = bot

	t.listenToChannels()
}
