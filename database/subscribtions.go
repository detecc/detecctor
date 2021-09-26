package database

import (
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// GetSubscribedChats get all the chats that include subscription(s) where the nodeId == nodeId and command == command
// or either node == * or command == *.
func GetSubscribedChats(nodeId, command string) ([]Chat, error) {
	return getChats(
		bson.M{"subscriptions.nodeId": bson.M{
			operator.In: bson.A{nodeId, "*"},
		},
			"subscriptions.command": bson.M{
				operator.In: bson.A{command, "*"},
			},
		},
	)
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

func SubscribeTo(chatId int64, nodes []string, commands []string) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	// check if there is a subscription to all nodes and commands
	if len(chat.Subscriptions) == 1 {
		firstSubscription := chat.Subscriptions[0]
		if firstSubscription.Node == "*" && firstSubscription.Command == "*" {
			// replace the all subscription with provided subscriptions
			chat.Subscriptions = createSubscriptions(nodes, commands)
			err = updateChat(chat)
			if err != nil {
				return err
			}
			return nil
		}
	}

	subs := createSubscriptions(nodes, commands)

	for _, sub := range subs {
		isDuplicateFound := false
		// check if there is an existing subscription for a node and command
		for _, subscription := range chat.Subscriptions {
			if sub.Node == subscription.Node && sub.Command == sub.Command {
				isDuplicateFound = true
			}
		}
		if !isDuplicateFound {
			chat.Subscriptions = append(chat.Subscriptions, sub)
		}
	}

	err = updateChat(chat)
	if err != nil {
		return err
	}

	return nil
}

func createSubscriptions(nodes []string, commands []string) []Subscription {
	var subscriptions []Subscription
	for _, nodeId := range nodes {
		// check if the node exists
		_, err := GetClientWithServiceNodeKey(nodeId)
		if err != nil {
			log.Println("Error creating a subscription for", nodeId, ":node doesnt exist")
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
