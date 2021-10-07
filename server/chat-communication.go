package server

import (
	"fmt"
	"github.com/detecc/detecctor/i18n"
	"github.com/detecc/detecctor/shared"
	"log"
	"reflect"
)

// SendMessageToChat is used to expose the server replyToChat method to the Plugins.
func SendMessageToChat(chatId int64, messageType int, content interface{}) error {
	switch messageType {
	case shared.TypeMessage, shared.TypePhoto, shared.TypeAudio:
		srv.replyToChat(chatId, content, messageType)
		return nil
	default:
		return fmt.Errorf("unsupported message type")
	}
}

func MakeTranslationMap(messageId string, plural interface{}, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	translatedMap := make(map[string]interface{})
	translatedMap["messageId"] = messageId
	translatedMap["data"] = data
	translatedMap["plural"] = plural
	return translatedMap
}

//replyToChat replies to telegram chat with a message.
func (s *server) replyToChat(chatId int64, content interface{}, contentType int) {

	if contentType == shared.TypeMessage && reflect.ValueOf(content).Kind() == reflect.Map {
		// if the content type is a message and contains a map, translate the message
		translatedMessage, err := i18n.TranslateReplyMessage(chatId, content)
		if err != nil {
			log.Println("Error translating the message", err)
			return
		}

		content = translatedMessage
	}

	s.replyChannel <- shared.Reply{
		ChatId:    chatId,
		ReplyType: contentType,
		Content:   content,
	}
}
