package bot

import (
	"fmt"
	. "github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/database"
	"log"
	"strings"
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

func (proxy *Proxy) GetReplyChannel() chan Reply {
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
			// ignore any non-Message Updates
			proxy.processMessage(message)
			break
		}
	}
}

// processMessage processes a message. Adds the necessary information to the database and sends the command to the TCP server.
// If an error occurs during processing, the proxy will notify the bot.
func (proxy *Proxy) processMessage(message ProxyMessage) {
	log.Println("Processing a message")

	err := database.AddChatIfDoesntExist(message.ChatId, message.Username)
	if err != nil {
		log.Println("Error adding a chat:", err)
	}

	// add a new message to database
	_, err = database.NewMessage(message.ChatId, message.MessageId, message.Message)
	if err != nil {
		log.Println("Error adding a message:", err)
		return
	}

	//update last message id
	err = database.UpdateLastMessageId(message.MessageId)
	if err != nil {
		log.Println("Error updating the lastMessageId:", err)
		return
	}

	command, err := parseCommand(message.Message, message.ChatId)
	if err != nil {
		log.Println("Error parsing the message:", err)
		//if the command is invalid, notify the user through telegram
		proxy.ReplyChannel <- NewReplyBuilder().TypeMessage().ForChat(message.ChatId).WithContent(fmt.Sprintf("%s is not a command", message)).Build()
		return
	}

	proxy.CommandsChannel <- command
}

// parseCommand parses the text as a command, where the command is structured as /command arg1 arg2 arg3.
// returns a Command struct containing the name of the command and the arguments provided : ["/command", "arg1", "arg2", "arg3"]
func parseCommand(text string, chatId string) (Command, error) {
	if !strings.HasPrefix(text, "/") {
		return Command{}, fmt.Errorf("not a command: %s", text)
	}
	args := strings.Split(text, " ")
	cmdBuilder := NewCommandBuilder()

	if len(args) == 1 {
		return cmdBuilder.WithName(args[0]).FromChat(chatId).Build(), nil
	}

	return cmdBuilder.WithName(args[0]).WithArgs(args[1:]).FromChat(chatId).Build(), nil
}
