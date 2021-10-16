package server

import (
	"fmt"
	"github.com/Allenxuxu/gev/connection"
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/i18n"
	"github.com/detecc/detecctor/shared"
	cache2 "github.com/patrickmn/go-cache"
	"log"
	"time"
)

// sendMessage send a message to a client.
func (s *server) sendMessage(message shared.Payload) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if message.ServiceNodeKey == "" {
		return fmt.Errorf("ServiceNodeKey is not set")
	}

	if message.Id == "" {
		return fmt.Errorf("payload id not set")
	}

	// find the target client for the message
	conn, err := s.getConnection(message.ServiceNodeKey)
	if err != nil {
		return err
	}

	encodedPayload, err := shared.EncodePayload(&message)
	if err != nil {
		return err
	}

	log.Println("Sending a message to a service node", message.ServiceNodeKey)
	return conn.Send([]byte(encodedPayload + "\n"))
}

func (s *server) sendToClients(chatId string, payloads ...shared.Payload) {
	if payloads != nil {
		for _, payload := range payloads {
			generatePayloadId(&payload, chatId)
			messageErr := s.sendMessage(payload)
			if messageErr != nil {
				log.Println("Could not send message to the client:", messageErr)
				translationMap := i18n.NewTranslationMap("UnableToSendMessage", i18n.AddData("ServiceNodeKey", payload.ServiceNodeKey), i18n.AddData("Error", messageErr.Error()))
				s.replyToChat(chatId, translationMap, api.TypeMessage)
			}
		}
	}
}

// handleMessage Handle a reply from the client
func (s *server) handleMessage(c *connection.Connection, payload shared.Payload) {
	clientId, _ := shared.ParseIP(c.PeerAddr())
	isNodeAuthorized := database.IsClientAuthorized(clientId)
	log.Println("Handling message from", clientId)

	switch payload.Command {
	case AuthCommand:
		if isNodeAuthorized {
			log.Println("Client is already authorized")
			payload.SetError(fmt.Errorf("already authorized"))

			//generate a dummy id
			generatePayloadId(&payload, "0")

			// reply to the client authorization request
			err := s.sendMessage(payload)
			if err != nil {
				log.Println(err)
			}
			return
		}

		err := s.authorizeClient(clientId, payload)
		if err != nil {
			log.Println("Client couldn't be authorized:", err)
			payload.SetError(err)
			// s.sendMessage(payload)
			return
		}
		break
	default:
		if !isNodeAuthorized {
			log.Println("Client is not authorized.")
			payload.SetError(fmt.Errorf("not authorized"))

			//generate a dummy id
			generatePayloadId(&payload, "0")

			err := s.sendMessage(payload)
			if err != nil {
				log.Println(err)
				return
			}
			c.ShutdownWrite()
			return
		}

		// if the payload id is empty, forward to all subscribed chats
		if payload.Id == "" {
			chatIds := s.getSubscribedChats(payload.ServiceNodeKey, payload.Command)
			if chatIds != nil {
				s.sendToSubscribedChats(chatIds, &payload)
			}
			return
		}

		// payloadId that is not empty usually means responding to a request from the user/chat
		chatId, isFound := cache.Memory().Get(payload.Id)
		if !isFound && payload.Id != "" {
			log.Print("chatId not found for payload id", payload.Id)
			return
		}

		s.sendToSubscribedChats([]string{chatId.(string)}, &payload)
		break
	}
}

// storeClient Remember a client connection for status updates
func (s *server) storeClient(conn *connection.Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.conn.PushBack(conn)
	conn.SetContext(e)

	clientId, _ := shared.ParseIP(conn.PeerAddr())
	cache.Memory().Set(clientId, conn, cache2.NoExpiration)

	// Check if the client already exists in the database
	log.Println("Adding the client to the database:", clientId)
	database.CreateClientIfNotExists(clientId, conn.PeerAddr(), "")

	err := database.UpdateClientStatus(clientId, database.StatusUnauthorized)
	if err != nil {
		log.Println("Cannot update the client status")
		conn.ShutdownWrite()
		return
	}
}

// generatePayloadId Generates a UUID for an outbound Payload to map the response to the ChatId
func generatePayloadId(payload *shared.Payload, chatId string) {
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
