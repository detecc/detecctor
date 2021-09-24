package server

import (
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"strings"
)

const (
	NodeId         = "nodes"
	Command        = "commands"
	NotifyInterval = "notifyInterval"
)

func getChatsSubscribed(nodeId, command string) []int64 {
	// todo get from cache?

	// fetch from database
	var chatIds []int64
	chats, err := database.GetChatsToSendTheNotificationTo(nodeId, command)
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

	// send the response to the Telegram Bot
	for _, chatId := range chatIds {
		s.replyToChat(chatId, pluginResponse.Content, pluginResponse.ReplyType)
	}
}

func (s *server) handleSubscription(command bot.Command) {
	var (
		nodes    = []string{"*"}
		commands = []string{"*"}
	)

	// /sub type value -> /sub command /get_status -> /sub node1,node2 /get_status,/get_cpu_usage
	settings := interpretSubscriptionCommand(command)

	if len(settings) == 0 {
		err := database.SubscribeToAll(command.ChatId)
		if err != nil {
			log.Println(err)
			s.replyToChat(command.ChatId, "an error occurred during subscription", shared.TypeMessage)
		}
		// subscribe to all
		s.replyToChat(command.ChatId, "subscribed to all nodes and commands", shared.TypeMessage)
		return
	}

	for key, value := range settings {
		switch key {
		case NodeId:
			nodes = strings.Split(value, ",")
			break
		case Command:
			commands = strings.Split(value, ",")
			break
		case NotifyInterval:
			break
		}
	}

	err := database.AddSubscriptions(command.ChatId, nodes, commands)
	if err != nil {
		s.replyToChat(command.ChatId, "something went wrong while subscribing", shared.TypeMessage)
	}
}

func (s *server) handleUnsubscription(command bot.Command) {
	var (
		nodes    = []string{"*"}
		commands = []string{"*"}
	)
	settings := interpretSubscriptionCommand(command)

	if len(settings) == 0 {
		// unsubscribe from all
		err := database.UnSubscribeFromAll(command.ChatId)
		if err != nil {
			log.Println("error unsubscribing from all nodes and commands:", err)
			s.replyToChat(command.ChatId, "could not unsubscribe from all nodes and commands", shared.TypeMessage)
		}
		s.replyToChat(command.ChatId, "successfully unsubscribed from all nodes and commands", shared.TypeMessage)
	}

	for key, value := range settings {
		switch key {
		case NodeId:
			nodes = strings.Split(value, ",")
			break
		case Command:
			commands = strings.Split(value, ",")
			break
		case NotifyInterval:
			break
		}
	}
	err := database.UnSubscribeFrom(command.ChatId, nodes, commands)
	if err != nil {
		s.replyToChat(command.ChatId, "could not unsubscribe from nodes and commands", shared.TypeMessage)
	}
}

func interpretSubscriptionCommand(command bot.Command) map[string]string {
	keyValues := map[string]string{}
	for _, args := range command.Args {
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
