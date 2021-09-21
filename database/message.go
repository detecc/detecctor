package database

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func getMessage(filter interface{}) (*Message, error) {
	msg := &Message{}
	messageCollection := mgm.Coll(msg)

	err := messageCollection.First(filter, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func addMessage(message *Message) error {
	msgCollection := &Message{}

	msg, err := GetMessageWithId(message.MessageId)
	if err == nil || msg != nil {
		return fmt.Errorf("duplicate message found")
	}

	return mgm.Coll(msgCollection).Create(message)
}

func GetMessageFromChat(chatId int) (*Message, error) {
	return getMessage(bson.M{"chatId": chatId})
}

func GetMessagesFromChat(chatId string) ([]Message, error) {
	var (
		msg               = &Message{}
		messageCollection = mgm.Coll(msg)
		results           []Message
	)

	// find all messages with the chatId
	cursor, err := messageCollection.Find(context.TODO(), bson.D{{"chatId", chatId}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	for _, result := range results {
		log.Println(result)
	}

	return results, nil
}

func GetMessageWithId(messageId int) (*Message, error) {
	return getMessage(bson.M{"messageId": messageId})
}

func NewMessage(chatId int, messageId int, content string) (*Message, error) {
	message := &Message{
		ChatId:    chatId,
		Content:   content,
		MessageId: messageId,
	}

	err := addMessage(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
