package database

import (
	"fmt"
	"github.com/detecc/detecctor/i18n"
	"golang.org/x/text/language"
)

// GetLanguage for a chat
func GetLanguage(chatId int64) (string, error) {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return "", err
	}

	return chat.Language, nil
}

// SetLanguage changes the language preference from a default one.
func SetLanguage(chatId int64, lang string) error {
	chat, err := GetChatWithId(chatId)
	if err != nil {
		return err
	}

	tag, _ := language.MatchStrings(i18n.Matcher, lang)
	if tag.String() == lang { // if the language is supported, update the chat
		chat.Language = tag.String()
		return updateChat(chat)
	}

	return fmt.Errorf("unsupported language")
}
