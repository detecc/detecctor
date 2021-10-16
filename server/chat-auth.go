package server

import (
	"crypto/rand"
	"fmt"
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/database"
	"log"
	"time"
)

// authChat Authorize the user (chat). If it is the first authentication attempt, generate a token and log it.
// The token is cached and expires in 5 minutes. If the user manages to provide the same token within a 5-minute time frame,
// it is considered that the user is authorized to use the server.
func (s *server) authChat(token string, chatId string) string {
	var returnMessage = "InvalidToken"
	cachedTokenId := fmt.Sprintf("auth-token-%s", chatId)

	if database.IsChatAuthorized(chatId) {
		return "AuthorizationError"
	}

	// check if the token is in the cache and if it matches the provided token
	cachedToken, isFound := cache.Memory().Get(cachedTokenId)
	if isFound && cachedToken.(string) == token {
		returnMessage = "ChatAuthorized"

		err := database.AuthorizeChat(chatId)
		if err != nil {
			log.Println(err)
			return "AuthorizationError"
		}

		cache.Memory().Delete(cachedTokenId)
		return returnMessage
	}

	if !isFound && token == "" {
		// generate a token
		GenerateChatAuthenticationToken(chatId)
		returnMessage = "GeneratedToken"
	}

	return returnMessage
}

// GenerateChatAuthenticationToken Generate an authorization token for a chat and log it. The token is cached and expires after 5 minutes.
func GenerateChatAuthenticationToken(chatId string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	authToken := fmt.Sprintf("%x", b)
	cache.Memory().Set(fmt.Sprintf("auth-token-%s", chatId), authToken, time.Minute*5)
	log.Println(authToken)
}
