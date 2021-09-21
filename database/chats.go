package database

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func getChat(filter bson.M) (*Chat, error) {
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

func AuthorizeChat(chatId int64) error {
	chat, err := getChat(bson.M{"chatId": chatId})
	if err != nil {
		return err
	}
	chat.IsAuthorized = true
	return updateChat(chat)
}

func RevokeChatAuthorization(chatId int64) error {
	chat, err := getChat(bson.M{"chatId": chatId})
	if err != nil {
		return err
	}
	chat.IsAuthorized = false
	return updateChat(chat)
}

func GetChatWithId(chatId int64) (*Chat, error) {
	return getChat(bson.M{"chatId": chatId})
}

func GetChats() ([]Chat, error) {
	var (
		chat    = &Chat{}
		results []Chat
	)

	// find all messages with the chatId
	cursor, err := mgm.Coll(chat).Find(context.TODO(), bson.M{})

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func AddChat(chatId int64, name string) error {
	if chatId < 0 {
		return nil
	}
	chat := &Chat{
		ChatId: chatId,
		Name:   name,
	}

	return mgm.Coll(&Chat{}).Create(chat)
}
