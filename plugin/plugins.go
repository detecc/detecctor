package plugin

import (
	"sn-bot/shared"
)

type Handler interface {
	// Response is called when the clients have responded and should
	// return a string to send as a reply to the bot
	Response(payload shared.Payload) shared.Reply
	// Execute method is called when the bot command matches GetCmdName's result.
	// The bot passes the string arguments to the method.
	// The execute method must return Payload array ready to be sent to the clients.
	Execute(args ...string) ([]shared.Payload, error)
}
