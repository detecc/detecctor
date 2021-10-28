package server

import (
	"fmt"
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/i18n"
	"log"
	"reflect"
)

// SendMessageToChat is used to expose the server replyToChat method to the Plugins.
// To translate the message, use the NewTranslationMap function to generate the content and use shared.TypeMessage as the message type.
func SendMessageToChat(chatId string, messageType int, content interface{}) error {
	switch messageType {
	case api.TypeMessage, api.TypePhoto, api.TypeAudio:
		srv.replyToChat(chatId, content, messageType)
		return nil
	default:
		return fmt.Errorf("unsupported message type")
	}
}

//replyToChat replies to telegram chat with a message. If the content is a map and the type is a message,
// the content/message will be translated prior to sending the message to the chat.
func (s *server) replyToChat(chatId string, content interface{}, contentType int) {

	if contentType == api.TypeMessage && reflect.ValueOf(content).Kind() == reflect.Struct {
		// if the content type is a message and contains a map, translate the message
		translatedMessage, err := TranslateReplyMessage(chatId, content)
		if err != nil {
			log.Println("Error translating the message", err)
			return
		}

		content = translatedMessage
	}

	s.replyChannel <- api.Reply{
		ChatId:    chatId,
		ReplyType: contentType,
		Content:   content,
	}
}

// TranslateReplyMessage remaps the content and translates the message using i18n.
// The translation is dependent on the chat language.
func TranslateReplyMessage(chatId string, content interface{}) (string, error) {
	translationMap := content.(i18n.TranslationMap)
	messageId := translationMap.MessageId
	data := translationMap.Data
	plural := translationMap.Plural

	// get the preferred language for the chat.
	lang, err2 := database.GetLanguage(chatId)
	if err2 != nil {
		return "", err2
	}

	// see if the translation is available
	localize, err := i18n.Localize(lang, messageId, data, plural)
	if err != nil {
		log.Println("Error localizing the translationMap", err)
		return "", err
	}

	return localize, nil
}
