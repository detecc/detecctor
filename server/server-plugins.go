package server

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/i18n"
	"github.com/detecc/detecctor/server/middleware"
	plugin2 "github.com/detecc/detecctor/server/plugin"
	"log"
	"strings"
)

// validateCommand check if the chat is authorized to perform a command.
func (s *server) validateCommand(command api.Command) error {
	if !database.IsChatAuthorized(command.ChatId) && command.Name != "/auth" {
		return fmt.Errorf("chat is not authorized")
	}

	return nil
}

// handleCommand handles the invocation of the Plugin.Execute method and sends the payloads produced to the designated clients.
func (s *server) handleCommand(command api.Command) {

	cmdErr := s.validateCommand(command)
	if cmdErr != nil {
		log.Println("Chat is not authorized to execute", command.Name)
		s.replyToChat(command.ChatId, i18n.NewTranslationMap("ChatUnauthorized"), api.TypeMessage)
		return
	}

	switch command.Name {
	case AuthCommand:
		var token = ""

		if len(command.Args) >= 1 {
			token = command.Args[0]
		}

		message := s.authChat(token, command.ChatId)
		s.replyToChat(command.ChatId, i18n.NewTranslationMap(message), api.TypeMessage)
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
				s.replyToChat(command.ChatId, fmt.Sprintf("An error occured while setting the language: %v.", err), api.TypeMessage)
				return
			}

			s.replyToChat(command.ChatId, fmt.Sprintf("Successfully set the language to: %s.", lang), api.TypeMessage)
			break
		}

		s.replyToChat(command.ChatId, i18n.NewTranslationMap("InvalidArguments"), api.TypeMessage)
		break
	default:
		s.executeCommand(command)
		break
	}
}

//executeMiddleware execute middleware registered to the plugin
func (s *server) executeMiddleware(ctx context.Context, metadata plugin2.Metadata) error {
	middlewareErr := middleware.GetMiddlewareManager().Chain(ctx, metadata.Middleware...)

	if middlewareErr != nil && !strings.Contains(middlewareErr.Error(), "not found") {
		log.Println(middlewareErr)
		return middlewareErr
	}
	return nil
}

//executeCommand executes the plugin associated with the command and sends a message to the client(s).
func (s *server) executeCommand(command api.Command) {
	//check if the plugin exists
	plugin, err := plugin2.GetPluginManager().GetPlugin(command.Name)
	if err != nil {
		log.Println("Plugin with command", command.Name, "doesnt exist")
		s.replyToChat(command.ChatId, i18n.NewTranslationMap("UnsupportedCommand", i18n.AddData("Command", command.Name)), api.TypeMessage)
		return
	}

	pluginMetadata := plugin.GetMetadata()

	ctx := context.WithValue(context.Background(), "", "")
	middlewareErr := s.executeMiddleware(ctx, pluginMetadata)

	// invoke the Plugin.Execute method
	payloads, err := plugin.Execute(command.Args...)

	if middlewareErr != nil || err != nil {
		var errorOpt i18n.TranslationOptions
		log.Println("plugin produced an error:", err)

		if err != nil {
			errorOpt = i18n.AddData("Error", err.Error())
		} else if middlewareErr != nil {
			errorOpt = i18n.AddData("Error", middlewareErr.Error())
		}
		translations := i18n.NewTranslationMap("PluginExecutionFailed", errorOpt)

		s.replyToChat(command.ChatId, translations, api.TypeMessage)
		return
	}

	switch pluginMetadata.Type {
	case plugin2.PluginTypeServerOnly:
		// do nothing. If the plugin needs to send something to the user, it should call SendMessageToChat method.
		break
	case plugin2.PluginTypeServerClient:
		// send the payloads to the clients
		s.sendToClients(command.ChatId, payloads...)
		break
	default:
		log.Println("Invalid plugin type:", pluginMetadata.Type)
		translationMap := i18n.NewTranslationMap("InvalidPluginType", i18n.AddData("Plugin", command.Name), i18n.AddData("PluginType", pluginMetadata.Type))

		s.replyToChat(command.ChatId, translationMap, api.TypeMessage)
		return
	}
}
