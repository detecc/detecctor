package server

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/server/middleware"
	plugin2 "github.com/detecc/detecctor/server/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"strings"
)

// handleCommand handles the invocation of the Plugin.Execute method and sends the payloads produced to the designated clients.
func (s *server) handleCommand(command bot.Command) {

	cmdErr := s.validateCommand(command)
	if cmdErr != nil {
		log.Println("Chat is not authorized to execute", command.Name)
		s.replyToChat(command.ChatId, MakeTranslationMap("ChatUnauthorized", nil, nil), shared.TypeMessage)
		return
	}

	switch command.Name {
	case AuthCommand:
		var token = ""

		if len(command.Args) >= 1 {
			token = command.Args[0]
		}

		message := s.authChat(token, command.ChatId)
		s.replyToChat(command.ChatId, MakeTranslationMap(message, nil, nil), shared.TypeMessage)
		break
	case SubscribeCommand, "/subscribe":
		s.handleSubscription(command)
		break
	case UnSubscribeCommand, "/unsubscribe":
		s.handleUnsubscription(command)
		break
	case LanguageCommand, LanguageCommandShort:
		if len(command.Args) >= 1 {
			lang := command.Args[0]

			err := database.SetLanguage(command.ChatId, lang)
			if err != nil {
				s.replyToChat(command.ChatId, fmt.Sprintf("An error occured while setting the language: %v.", err), shared.TypeMessage)
				return
			}

			s.replyToChat(command.ChatId, fmt.Sprintf("Successfully set the language to: %s.", lang), shared.TypeMessage)
			break
		}

		s.replyToChat(command.ChatId, MakeTranslationMap("InvalidArguments", nil, nil), shared.TypeMessage)
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
		s.replyToChat(chatId, "an error occurred while executing middleware:", shared.TypeMessage)
		return
	}
}

//executePlugin executes the plugin associated with the command and sends a message to the client(s).
func (s *server) executePlugin(command bot.Command) {
	//check if the plugin exists
	plugin, err := plugin2.GetPluginManager().GetPlugin(command.Name)
	if err != nil {
		log.Println("Plugin with command", command.Name, "doesnt exist")

		dataMap := make(map[string]interface{})
		dataMap["Command"] = command.Name

		s.replyToChat(command.ChatId, MakeTranslationMap("UnsupportedCommand", nil, dataMap), shared.TypeMessage)
		return
	}

	pluginMetadata := plugin.GetMetadata()
	ctx := context.TODO()
	s.executeMiddleware(command.ChatId, ctx, pluginMetadata)

	// invoke the Plugin.Execute method
	payloads, err := plugin.Execute(command.Args...)
	if err != nil {
		log.Println("plugin produced an error:", err)
		dataMap := make(map[string]interface{})
		dataMap["Error"] = err.Error()

		s.replyToChat(command.ChatId, MakeTranslationMap("PluginExecutionFailed", nil, dataMap), shared.TypeMessage)
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
					log.Println("Could not send message to the client:", err)

					dataMap := make(map[string]interface{})
					dataMap["ServiceNodeKey"] = payload.ServiceNodeKey
					dataMap["Error"] = err.Error()

					s.replyToChat(command.ChatId, MakeTranslationMap("UnableToSendMessage", nil, dataMap), shared.TypeMessage)
				}
			}
		}
		break
	default:
		log.Println("Invalid plugin type:", pluginMetadata.Type)

		dataMap := make(map[string]interface{})
		dataMap["Plugin"] = command.Name
		dataMap["PluginType"] = pluginMetadata.Type

		s.replyToChat(command.ChatId, MakeTranslationMap("InvalidPluginType", nil, dataMap), shared.TypeMessage)
		return
	}
}
