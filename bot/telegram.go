package bot

import (
	"fmt"
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/database"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
)

// Telegram is a wrapper for the Telegram bot API.
type Telegram struct {
	Token          string
	botAPI         *telegram.BotAPI
	messageChannel chan api.ProxyMessage
}

// Start listening to the bot updates and the updates from the TCP server
func (t *Telegram) Start() {
	bot, err := telegram.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	t.botAPI = bot
	t.messageChannel = make(chan api.ProxyMessage)
}

func (t *Telegram) GetMessageChannel() chan api.ProxyMessage {
	return t.messageChannel
}

// ListenToChannels listens for incoming data from replyChannel and from telegram bot messages
func (t *Telegram) ListenToChannels() {
	log.Printf("Authorized on account %s", t.botAPI.Self.UserName)
	message, err := database.GetStatistics()

	lastMessageId := 0
	if err == nil {
		messageId, err := strconv.Atoi(message.UpdateId)
		if err == nil {
			lastMessageId = messageId
		} else {
			log.Println(err)
		}
	}

	u := telegram.NewUpdate(lastMessageId)
	u.Timeout = 60

	updates, err := t.botAPI.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case update := <-updates:
			if update.Message == nil || update.Message.Entities == nil || len(*update.Message.Entities) == 0 {
				return
			}

			for _, entity := range *update.Message.Entities {
				if entity.Type == "bot_command" {
					chatId := fmt.Sprintf("%d", update.Message.Chat.ID)
					messageId := fmt.Sprintf("%d", update.Message.MessageID)
					builder := api.NewMessageBuilder().WithId(chatId).FromUser(update.Message.Chat.UserName).WithMessage(messageId, update.Message.Text)
					t.messageChannel <- builder.Build()
				}
			}

			break
		}
	}
}

func (t *Telegram) ReplyToChat(replyMessage api.Reply) {
	var msg telegram.Chattable

	chatId, err := strconv.Atoi(replyMessage.ChatId)
	if err != nil {
		log.Println(err)
		return
	}

	switch replyMessage.ReplyType {
	case api.TypeMessage:
		if replyMessage.Content != nil {
			msg = telegram.NewMessage(int64(chatId), replyMessage.Content.(string))
		}
		break
	case api.TypePhoto:
		msg = telegram.NewPhotoUpload(int64(chatId), replyMessage.Content)
		break
	default:
		return
	}

	if msg != nil {
		log.Println("Replying to chat", replyMessage.ChatId)
		_, err := t.botAPI.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
