package server

import (
	"fmt"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/server/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"strings"
)

const (
	NodeId         = "nodes"
	Command        = "commands"
	NotifyInterval = "notifyInterval"
)

// getSubscribedChats get all the chats that are subscribed to the nodeId, command combo.
func (s *server) getSubscribedChats(nodeId, command string) []int64 {
	// todo get from cache?

	// fetch from database
	var chatIds []int64
	chats, err := database.GetSubscribedChats(nodeId, command)
	if err != nil {
		log.Println(err)
		return chatIds
	}

	for _, chat := range chats {
		chatIds = append(chatIds, chat.ChatId)
	}
	return chatIds
}

// sendToSubscribedChats sends the response of the plugin to each chat that is subscribed to the command or node.
func (s *server) sendToSubscribedChats(chatIds []int64, payload *shared.Payload) {
	// send a notification to the users about the failure
	if payload.Success == false {
		for _, chatId := range chatIds {
			s.replyToChat(chatId, payload.Error, shared.TypeMessage)
		}
		return
	}

	// Gets the response from the plugin
	mPlugin, err := plugin.GetPluginManager().GetPlugin(payload.Command)
	if err != nil {
		log.Println("Plugin doesnt exist")
		return
	}
	pluginResponse := mPlugin.Response(*payload)

	// send the response to the Bot
	for _, chatId := range chatIds {
		s.replyToChat(chatId, pluginResponse.Content, pluginResponse.ReplyType)
	}
}

// handleSubscription handles the subscribe command and its arguments. Updates the Chat and notifies the user of the result.
func (s *server) handleSubscription(command bot.Command) {
	settings := interpretSubscriptionCommand(command.Args)

	// subscribe to all
	if len(settings) == 0 {
		err := database.SubscribeToAll(command.ChatId)
		if err != nil {
			log.Println("An error occurred during subscription;", err)
			data := make(map[string]interface{})
			data["Error"] = err.Error()

			s.replyToChat(command.ChatId, MakeTranslationMap("SubscriptionFail", nil, data), shared.TypeMessage)
			return
		}

		s.replyToChat(command.ChatId, MakeTranslationMap("SubscriptionSuccess", nil, nil), shared.TypeMessage)
		return
	}

	nodes, commands := getNodesAndCommands(settings)
	err := database.SubscribeTo(command.ChatId, nodes, commands)
	if err != nil {
		log.Println("An error occurred during subscription;", err)
		data := make(map[string]interface{})
		data["Error"] = err.Error()
		s.replyToChat(command.ChatId, MakeTranslationMap("SubscriptionFail", nil, data), shared.TypeMessage)
		return
	}
	s.replyToChat(command.ChatId, MakeTranslationMap("SubscriptionSuccess", nil, nil), shared.TypeMessage)
}

// handleUnsubscription handles unsubscribe command.
func (s *server) handleUnsubscription(command bot.Command) {
	settings := interpretSubscriptionCommand(command.Args)

	// if there are no arguments, unsubscribe from all
	if len(settings) == 0 {
		err := database.UnSubscribeFromAll(command.ChatId)
		if err != nil {
			log.Println("Error unsubscribing from all nodes and commands:", err)
			data := make(map[string]interface{})
			data["Error"] = err.Error()
			s.replyToChat(command.ChatId, MakeTranslationMap("UnsubscribeFail", nil, data), shared.TypeMessage)
			return
		}
		s.replyToChat(command.ChatId, MakeTranslationMap("UnsubscribeSuccess", nil, nil), shared.TypeMessage)
		return
	}

	nodes, commands := getNodesAndCommands(settings)
	err := database.UnSubscribeFrom(command.ChatId, nodes, commands)
	if err != nil {
		log.Println("error unsubscribing from nodes and commands:", err)
		data := make(map[string]interface{})
		data["Error"] = err.Error()
		s.replyToChat(command.ChatId, MakeTranslationMap("UnsubscribeFail", nil, data), shared.TypeMessage)
		return
	}
	s.replyToChat(command.ChatId, MakeTranslationMap("UnsubscribeSuccess", nil, nil), shared.TypeMessage)
}

// interpretSubscriptionCommand interpret the arguments of the command and return key-value pairs to further processing.
// Example command: /sub nodes=node1,node2,node3 commands=/auth,/get_status, where the equals sign means key-value mapping.
// It will remove any excess or unsupported arguments.
func interpretSubscriptionCommand(args []string) map[string]string {
	keyValues := map[string]string{}
	for _, args := range args {
		keyValue := strings.Split(args, "=")
		if len(keyValue) >= 2 {
			switch keyValue[0] {
			case NodeId, NotifyInterval, Command:
				keyValues[keyValue[0]] = keyValue[1]
				break
			}
		}
	}
	return keyValues
}

// getNodesAndCommands gets all the nodes and command entries.
// Example command: /sub nodes=node1,node2,node3 commands=/auth,/get_status.
// It will return a list of nodes and commands.
func getNodesAndCommands(settings map[string]string) ([]string, []string) {
	var (
		nodes    = []string{"*"}
		commands = []string{"*"}
	)
	for key, value := range settings {
		switch key {
		case NodeId:
			nodes = strings.Split(value, ",")
			break
		case Command:
			commands = strings.Split(value, ",")
			for i, command := range commands {
				if !strings.HasPrefix(command, "/") {
					commands[i] = fmt.Sprintf("/%s", command)
				}
			}
			break
		case NotifyInterval:
			break
		}
	}
	return nodes, commands
}
