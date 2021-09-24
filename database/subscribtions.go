package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// GetChatsToSendTheNotificationTo get all the chats that include subscription(s):
// node == NodeId and command == command or nodeId == * and command == *.
// or node = NodeId and command = * or nodeId ==* and command == command.
func GetChatsToSendTheNotificationTo(nodeId, command string) ([]Chat, error) {
	return getChats(bson.M{"subscription": bson.E{}})
}

func SubscribeToAll(chatId int64) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	// this overwrites any previous subscriptions.
	chat.Subscriptions = []Subscription{
		{
			Node:    "*",
			Command: "*",
		},
	}

	err = updateChat(chat)
	if err != nil {
		return err
	}

	return nil
}

func AddSubscriptions(chatId int64, nodes []string, commands []string) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	// check if there is a subscription to everything
	if len(chat.Subscriptions) == 1 {
		firstSubscription := chat.Subscriptions[0]
		if firstSubscription.Node == "*" && firstSubscription.Command == "*" {
			// replace subscription all
			chat.Subscriptions = createSubscriptions(nodes, commands)
			err = updateChat(chat)
			if err != nil {
				return err
			}
			return nil
		} else {

		}
	}

	subs := createSubscriptions(nodes, commands)

	for _, sub := range subs {
		isDuplicateFound := false
		// check if there is a subscription for a node
		for _, subscription := range chat.Subscriptions {
			if sub.Node == subscription.Node && sub.Command == sub.Command {
				isDuplicateFound = true
			}
		}
		if !isDuplicateFound {
			chat.Subscriptions = append(chat.Subscriptions, sub)
		}

	}

	chat.Subscriptions = subs

	err = updateChat(chat)
	if err != nil {
		return err
	}

	return nil
}

func createSubscriptions(nodes []string, commands []string) []Subscription {
	var subscriptions []Subscription
	for _, nodeId := range nodes {
		_, err := GetClientWithServiceNodeKey(nodeId)
		// check if the node exists
		if err != nil {
			log.Println("node doesnt exist", err)
			continue
		}

		for _, command := range commands {
			subscriptions = append(subscriptions, Subscription{Node: nodeId, Command: command})
		}

	}

	return subscriptions
}

func UnSubscribeFromAll(chatId int64) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	// this overwrites any previous subscriptions.
	chat.Subscriptions = []Subscription{}

	err = updateChat(chat)
	if err != nil {
		return err
	}

	return nil
}

func UnSubscribeFrom(chatId int64, nodes []string, commands []string) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	if len(chat.Subscriptions) == 1 {
		firstSubscription := chat.Subscriptions[0]
		if firstSubscription.Node == "*" && firstSubscription.Command == "*" {
			return UnSubscribeFromAll(chatId)
		}
	}
	for i, subscription := range chat.Subscriptions {
		for _, node := range nodes {
			if node == "*" {
				chat.Subscriptions = append(chat.Subscriptions[:i], chat.Subscriptions[i+1:]...)
				continue
			}
			for _, command := range commands {
				if command == "*" || (command == subscription.Command && node == subscription.Node) {
					chat.Subscriptions = append(chat.Subscriptions[:i], chat.Subscriptions[i+1:]...)
				}
			}
		}
	}

	return nil
}
