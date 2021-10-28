package database

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func getChat(filter interface{}) (*Chat, error) {
	chat := &Chat{}
	chatCollection := mgm.Coll(chat)

	// Get the first doc of a collection using a filter
	err := chatCollection.First(filter, chat)
	if err != nil {
		log.Println("Error querying chats:", err)
		return nil, err
	}
	return chat, nil
}

func updateChat(chat *Chat) error {
	chatCollection := mgm.Coll(&Chat{})
	return chatCollection.Update(chat)
}

func IsChatAuthorized(chatId string) bool {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		log.Println("Error authenticating the chat:", err)
		return false
	}

	return chat.IsAuthorized
}

func AuthorizeChat(chatId string) error {
	chat, err := getChat(bson.M{"chatId": chatId})
	if err != nil {
		return err
	}
	chat.IsAuthorized = true
	return updateChat(chat)
}

func RevokeChatAuthorization(chatId string) error {
	chat, err := getChat(bson.M{"chatId": chatId})
	if err != nil {
		return err
	}
	chat.IsAuthorized = false
	return updateChat(chat)
}

func GetChatWithId(chatId string) (*Chat, error) {
	return getChat(bson.M{"chatId": chatId})
}

func GetChats() ([]Chat, error) {
	return getChats(bson.M{})
}

func getChats(filter interface{}) ([]Chat, error) {
	var (
		chat    = &Chat{}
		results []Chat
	)
	cursor, err := mgm.Coll(chat).Find(mgm.Ctx(), filter)
	if err = cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func AddChat(chatId string, name string) error {
	if chatId == "" {
		return nil
	}
	chat := &Chat{
		ChatId:        chatId,
		Name:          name,
		IsAuthorized:  false,
		Language:      "en",
		Subscriptions: []Subscription{},
	}

	return mgm.Coll(&Chat{}).Create(chat)
}

func AddChatIfDoesntExist(chatId string, name string) error {
	_, err := GetChatWithId(chatId)
	if err == nil {
		return fmt.Errorf("chat already exists")
	}

	return AddChat(chatId, name)
}
