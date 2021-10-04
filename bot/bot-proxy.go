package bot

import (
	"fmt"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/shared"
	"log"
)

// Start the bot and listening for messages from the TCP server and the bot.
func (proxy *Proxy) Start() {
	proxy.bot.Start()
	go proxy.listenForMessages()
	proxy.bot.ListenToChannels()
}

func (proxy *Proxy) GetCommandsChannel() chan Command {
	return proxy.CommandsChannel
}

func (proxy *Proxy) GetReplyChannel() chan shared.Reply {
	return proxy.ReplyChannel
}

// listenForMessages from the TCP server and the Bot.
// Incoming messages, from the TCP server are sent to the bot for further processing.
// Incoming messages from the Bot are processed and forwarded to the TCP server.
func (proxy *Proxy) listenForMessages() {
	for {
		select {
		case replyMessage := <-proxy.ReplyChannel:
			log.Println("Got reply:", replyMessage)
			proxy.bot.ReplyToChat(replyMessage)
			break
		case message := <-proxy.bot.GetMessageChannel():
			if message != nil {
				// ignore any non-Message Updates
				proxy.processMessage(
					message["chatId"].(int64),
					message["messageId"].(int),
					message["username"].(string),
					message["message"].(string),
				)
			}
			break
		}
	}

}

// processMessage processes a message. Adds the necessary information to the database and sends the command to the TCP server.
// If an error occurs during processing, the proxy will notify the bot.
func (proxy *Proxy) processMessage(chatId int64, messageId int, userName, message string) {
	log.Println("Processing a message")
	_, err := database.GetChatWithId(chatId)
	if err != nil {
		err := database.AddChat(chatId, userName)
		if err != nil {
			log.Println("Error adding a chat:", err)
			return
		}
	}

	// add a new message to database
	_, err = database.NewMessage(int(chatId), messageId, message)
	if err != nil {
		log.Println("Error adding a message:", err)
		return
	}

	//update last message id
	err = database.UpdateLastMessageId(messageId)
	if err != nil {
		log.Println("Error updating the lastMessageId:", err)
		return
	}

	command, err := parseCommand(message, chatId)
	if err != nil {
		log.Println("Error parsing the message:", err)
		//if the command is invalid, notify the user through telegram
		proxy.ReplyChannel <- shared.Reply{
			ChatId:    chatId,
			ReplyType: shared.TypeMessage,
			Content:   fmt.Sprintf("%s is not a command", message),
		}
		return
	}

	proxy.CommandsChannel <- command
}
