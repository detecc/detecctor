package bot

import (
	"github.com/Allenxuxu/gev/log"
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/config"
	"sync"
)

const (
	TelegramBot = "telegram"
)

var once = sync.Once{}
var proxy *Proxy

type (

	// Bot represents a chatbot (e.g. Telegram, Discord, Slack bot).
	Bot interface {
		// Start should initialize the bot to listen to the chat.
		Start()
		// ListenToChannels is called after the start function. It monitors the chat and should send the message data to the Message channel.
		ListenToChannels()
		// ReplyToChat after executing the command or if an error occurs.
		ReplyToChat(replyMessage api.Reply)
		// GetMessageChannel returns the channel the bot is supposed to notify the proxy of the incoming message.
		GetMessageChannel() chan api.ProxyMessage
	}

	// Proxy acts as a proxy between the bot implementation and the TCP server. It stores all the necessary data and listens to
	// both the bot and the TCP server.
	Proxy struct {
		bot             Bot
		CommandsChannel chan api.Command
		ReplyChannel    chan api.Reply
	}
)

func GetProxy(bot Bot) *Proxy {
	once.Do(func() {
		proxy = &Proxy{
			bot:             bot,
			CommandsChannel: make(chan api.Command),
			ReplyChannel:    make(chan api.Reply),
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
