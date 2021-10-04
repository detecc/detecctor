package bot

import (
	"fmt"
	"github.com/Allenxuxu/gev/log"
	"github.com/detecc/detecctor/config"
	"github.com/detecc/detecctor/shared"
	"strings"
	"sync"
)

const (
	TelegramBot = "telegram"
)

var once = sync.Once{}
var proxy *Proxy

type (

	// Bot represent a chatbot (e.g. Telegram, Discord, Slack bot).
	Bot interface {
		// Start should initialize the bot to listen to the chat.
		Start()
		// ListenToChannels is called after the start function. It monitors the chat and should send the message data to the Message channel.
		ListenToChannels()
		// ReplyToChat after executing the command or if an error occurs.
		ReplyToChat(replyMessage shared.Reply)
		// GetMessageChannel returns the channel the bot is supposed to notify the proxy of the incoming message.
		GetMessageChannel() chan map[string]interface{}
	}

	// Proxy acts as a proxy between the bot implementation and the TCP server. It stores all the necessary data and listens to
	// both the bot and the TCP server.
	Proxy struct {
		bot             Bot
		CommandsChannel chan Command
		ReplyChannel    chan shared.Reply
	}

	// Command consists of a Name and Args.
	// Example of a command: "/get_status serviceNode1 serviceNode2".
	// The command name is "/get_status", the arguments are ["serviceNode1", "serviceNode2"].
	Command struct {
		Name string
		// Args contains the arguments extracted from the Message sent through Telegram.
		Args   []string
		ChatId int64
	}
)

func GetProxy(bot Bot) *Proxy {
	once.Do(func() {
		proxy = &Proxy{
			bot:             bot,
			CommandsChannel: make(chan Command),
			ReplyChannel:    make(chan shared.Reply),
		}
	})
	return proxy
}

// NewBot create a new telegram bot
func NewBot(botConfiguration config.BotConfiguration) Bot {
	switch botConfiguration.Type {
	case TelegramBot:
		return &Telegram{
			Token: botConfiguration.Token,
		}
	default:
		log.Fatal("Unsupported bot type")
		return nil
	}
}

// parseCommand parses the text as a command, where the command is structured as /command arg1 arg2 arg3.
// returns a Command struct containing the name of the command and the arguments provided : ["/command", "arg1", "arg2", "arg3"]
func parseCommand(text string, chatId int64) (Command, error) {
	if !strings.HasPrefix(text, "/") {
		return Command{}, fmt.Errorf("not a command: %s", text)
	}
	args := strings.Split(text, " ")

	if len(args) == 1 {
		return Command{
			Name:   args[0],
			Args:   []string{},
			ChatId: chatId,
		}, nil
	}

	return Command{
		Name:   args[0],
		Args:   args[1:],
		ChatId: chatId,
	}, nil
}
