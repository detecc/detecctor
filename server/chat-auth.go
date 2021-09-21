package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/shared"
	"time"
)

//isChatAuthorized Check if the user (chat) is authorized in the database.
func (s *server) isChatAuthorized(chatId int64) bool {
	chat, err := database.GetChatWithId(chatId)
	if err != nil {
		log.Println("Error authenticating the chat:", err)
		return false
	}

	return chat.IsAuthorized
}

// authChat Authorize the user (chat). If it is the first authentication attempt, generate a token and log it.
// The token is cached and expires in 5 minutes. If the user manages to provide the same token within a 5-minute time frame,
// it is considered that the user is authorized to use the server.
func (s *server) authChat(token string, chatId int64) string {
	var returnMessage = "Invalid or expired token."
	// check if the token is in the cache and if it matches the provided token
	cachedToken, isFound := cache.Memory().Get(fmt.Sprintf("auth-token-%d", chatId))
	if isFound && cachedToken.(string) == token {
		returnMessage = "You have been authorized successfully!"
		err := database.AuthorizeChat(chatId)
		if err != nil {
			log.Println(err)
			returnMessage = "Error authorizing the chat."
		}

		return returnMessage
	}

	if !isFound && token == "" {
		// generate a token
		GenerateChatAuthenticationToken(chatId)
		returnMessage = "Generated a token for authentication."
	}

	return returnMessage
}

// GenerateChatAuthenticationToken Generate an authorization token for a chat and log it. The token is cached and expires after 5 minutes.
func GenerateChatAuthenticationToken(chatId int64) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	authToken := fmt.Sprintf("%x", b)
	cache.Memory().Set(fmt.Sprintf("auth-token-%d", chatId), authToken, time.Minute*5)
	log.Println(authToken)
}

// generatePayloadId Generates a UUID for an outbound Payload to map the response to the ChatId
func generatePayloadId(payload *shared.Payload, chatId int64) {
	//create a unique id for every server message
	uuid := shared.Uuid()
	log.Println("UUID:", uuid)
	if uuid == "" {
		// bad
		log.Println("uuid is empty")
		return
	}
	payload.Id = uuid
	//set the payload ID to chatId mapping
	cache.Memory().Set(uuid, chatId, time.Minute*5)
}
