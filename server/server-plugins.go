package server

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/middleware"
	plugin2 "github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"strings"
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

// handleCommand handles the invocation of the Plugin.Execute method and sends the payloads produced to the designated clients.
func (s *server) handleCommand(command bot.Command) {

	cmdErr := s.validateCommand(command)
	if cmdErr != nil {
		s.replyToChat(command.ChatId, "You are not authorized to send this command.", shared.TypeMessage)
		return
	}

	switch command.Name {
	case AuthCommand:
		var token = ""

		if len(command.Args) >= 1 {
			token = command.Args[0]
		}

		message := s.authChat(token, command.ChatId)
		s.replyToChat(command.ChatId, message, shared.TypeMessage)
		break
	case SubscribeCommand, "/subscribe":
		s.handleSubscription(command)
		break
	case UnSubscribeCommand, "/unsubscribe":
		s.handleUnsubscription(command)
		break
	default:
		s.executePlugin(command)
		break
	}
}

//executeMiddleware execute middleware registered to the plugin
func (s *server) executeMiddleware(chatId int64, ctx context.Context, metadata plugin2.Metadata) {
	middlewareErr := middleware.GetMiddlewareManager().Chain(ctx, metadata.Middleware...)
	if middlewareErr != nil && !strings.Contains(middlewareErr.Error(), "not found") {
		log.Println(middlewareErr)
		s.replyToChat(chatId, "an error occurred while executing", shared.TypeMessage)
		return
	}
}

//executePlugin executes the plugin associated with the command and sends a message to the client(s).
func (s *server) executePlugin(command bot.Command) {
	//check if the plugin exists
	plugin, err := plugin2.GetPluginManager().GetPlugin(command.Name)
	if err != nil {
		log.Println("Plugin with command", command.Name, "doesnt exist")
		s.replyToChat(command.ChatId, fmt.Sprintf("Command %s is unsupported.", command.Name), shared.TypeMessage)
		return
	}

	pluginMetadata := plugin.GetMetadata()
	ctx := context.TODO()
	s.executeMiddleware(command.ChatId, ctx, pluginMetadata)

	// invoke the Plugin.Execute method
	payloads, err := plugin.Execute(command.Args...)
	if err != nil {
		log.Println("plugin produced an error:", err)
		s.replyToChat(command.ChatId, "command could not be executed.", shared.TypeMessage)
		return
	}

	switch pluginMetadata.Type {

	case plugin2.PluginTypeServerOnly:
		// do nothing. If the plugin needs to send something to the user, it should call SendMessageToChat method.
		break
	case plugin2.PluginTypeServerClient:
		// send the payloads to the clients
		if payloads != nil {
			for _, payload := range payloads {
				generatePayloadId(&payload, command.ChatId)
				messageErr := s.sendMessage(payload)
				if messageErr != nil {
					couldNotSendMessage := fmt.Sprintf("could  not send message to %s: %v", payload.ServiceNodeKey, messageErr)
					log.Println(couldNotSendMessage)
					s.replyToChat(command.ChatId, couldNotSendMessage, shared.TypeMessage)
				}
			}
		}
		break
	default:
		log.Println("Invalid plugin type:", pluginMetadata.Type)
		s.replyToChat(command.ChatId, fmt.Sprintf("plugin %s unable to execute (invalid type).", command.Name), shared.TypeMessage)
		return
	}
}
