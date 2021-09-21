package bot

import (
	"errors"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/shared"
	"strings"
)

// listenToChannels listens for incoming data from replyChannel and from telegram bot messages
func (t *Telegram) listenToChannels() {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)
	message, err := database.GetStatistics()
	lastMessageId := 0
	if err == nil {
		lastMessageId = message.UpdateId
	}

	u := telegram.NewUpdate(lastMessageId)
	u.Timeout = 60

	updates, err := t.BotAPI.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case replyMessage := <-t.ReplyChannel:
			log.Println("Got reply:", replyMessage)
			t.replyToChat(replyMessage)
			break
		case update := <-updates:
			if update.Message == nil {
				return
			}

			// ignore any non-Message Updates
			t.ProcessMessage(update)
			break
		}
	}
}

func (t *Telegram) replyToChat(replyMessage shared.Reply) {
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
		_, err := t.BotAPI.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

// ProcessMessage process the telegram message and add the message to the database
func (t *Telegram) ProcessMessage(update telegram.Update) {
	message := update.Message
	chatId := message.Chat.ID
	if len(*message.Entities) == 0 {
		return
	}

	log.Println("Processing a message from telegram")
	_, err := database.GetChatWithId(chatId)
	if err != nil {
		err := database.AddChat(chatId, message.Chat.UserName)
		if err != nil {
			log.Println("Error adding a chat:", err)
			return
		}
	}
	// add a new message to database
	_, err = database.NewMessage(int(chatId), message.MessageID, message.Text)
	if err != nil {
		log.Println("Error adding a message:", err)
		return
	}
	//update last message id
	err = database.UpdateLastMessageId(update.UpdateID)
	if err != nil {
		log.Println("Error updating the lastMessageId:", err)
		return
	}

	for _, entity := range *message.Entities {
		command, err := t.parseCommand(message.Text, chatId)
		if err != nil {
			log.Println("telegram:", err)
			//if the command is invalid, notify the user through telegram
			t.ReplyChannel <- shared.Reply{
				ChatId:    chatId,
				ReplyType: shared.TypeMessage,
				Content:   fmt.Sprintf("%s is not a command", message.Text),
			}
			continue
		}
		if entity.Type == "bot_command" {
			t.CommandsChannel <- command
		}
	}
}

// parseCommand parses the text as a command, where the command is structured as /command arg1 arg2 arg3.
// returns a Command struct containing the name of the command and the arguments provided : ["/command", "arg1", "arg2", "arg3"]
func (t *Telegram) parseCommand(text string, chatId int64) (Command, error) {
	if !strings.HasPrefix(text, "/") {
		return Command{}, errors.New("not a command: " + text)
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
