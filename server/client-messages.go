package server

import (
	"github.com/Allenxuxu/gev/connection"
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
)

// handleMessage Handle a reply from the client
func (s *server) handleMessage(c *connection.Connection, payload shared.Payload) {
	clientId, _ := shared.ParseIP(c.PeerAddr())
	isNodeAuthorized := database.IsClientAuthorized(clientId)
	log.Println("handling message from", clientId)

	switch payload.Command {
	case "/auth":
		if isNodeAuthorized {
			log.Println("client is already authorized")
			payload.Error = "already authorized"
			payload.Success = false
			// reply to the client authorization request
			err := s.sendMessage(payload)
			if err != nil {
				log.Println(err)
				return
			}
			//c.ShutdownWrite()
			return
		}

		err := s.authorizeClient(clientId, payload)
		if err != nil {
			log.Println("client couldn't be authorized:", err)
			payload.Error = err.Error()
			payload.Success = false
			// s.sendMessage(payload)
			return
		}
		break
	default:
		if !isNodeAuthorized {
			log.Println("client is not authorized")
			payload.Error = "not authorized"
			payload.Success = false
			err := s.sendMessage(payload)
			if err != nil {
				log.Println(err)
				return
			}
			c.ShutdownWrite()
			return
		}

		chatId, isFound := cache.Memory().Get(payload.Id)
		if !isFound {
			// oopsie
			log.Print("chatId not found for payload id", payload.Id)
			return
		}

		// send a notification to the user about the failure
		if payload.Success == false {
			database.NewCommandResponse(payload.Id, nil, payload.Error)
			s.replyToChat(chatId.(int64), payload.Error, shared.TypeMessage)
			return
		}

		mPlugin, err := plugin.GetPluginManager().GetPlugin(payload.Command)
		if err != nil {
			database.NewCommandResponse(payload.Id, nil, err, payload.Error)
			log.Println("Plugin doesnt exist")
			return
		}
		// send the response to the Telegram Bot
		pluginResponse := mPlugin.Response(payload)
		database.NewCommandResponse(payload.Id, pluginResponse.Content, err, payload.Error)
		s.replyChannel <- pluginResponse
		break
	}
}
