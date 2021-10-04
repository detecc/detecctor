package bot

import (
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/shared"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// Telegram is a wrapper for the Telegram bot API.
type Telegram struct {
	Token          string
	botAPI         *telegram.BotAPI
	messageChannel chan map[string]interface{}
}

// Start listening to the bot updates and the updates from the TCP server
func (t *Telegram) Start() {
	bot, err := telegram.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	t.botAPI = bot
	t.messageChannel = make(chan map[string]interface{})
}

func (t *Telegram) GetMessageChannel() chan map[string]interface{} {
	return t.messageChannel
}

// ListenToChannels listens for incoming data from replyChannel and from telegram bot messages
func (t *Telegram) ListenToChannels() {
	log.Printf("Authorized on account %s", t.botAPI.Self.UserName)
	message, err := database.GetStatistics()
	lastMessageId := 0
	if err == nil {
		lastMessageId = message.UpdateId
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
					messageMap := make(map[string]interface{})
					messageMap["chatId"] = update.Message.Chat.ID
					messageMap["username"] = update.Message.Chat.UserName
					messageMap["message"] = update.Message.Text
					messageMap["messageId"] = update.Message.MessageID
					t.messageChannel <- messageMap
				}
			}

			break
		}
	}
}

func (t *Telegram) ReplyToChat(replyMessage shared.Reply) {
	var msg telegram.Chattable

	switch replyMessage.ReplyType {
	case shared.TypeMessage:
		if replyMessage.Content != nil {
			msg = telegram.NewMessage(replyMessage.ChatId, replyMessage.Content.(string))
		} else {
			msg = telegram.NewMessage(replyMessage.ChatId, "")
		}
		break
	case shared.TypePhoto:
		msg = telegram.NewPhotoUpload(replyMessage.ChatId, replyMessage.Content)
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
