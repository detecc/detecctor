package plugin

import (
	"github.com/detecc/detecctor/shared"
)

// constants for Metadata.Type
const (
	PluginTypeServerOnly   = "serverOnly"
	PluginTypeServerClient = "serverClient"
)

type (
	Handler interface {

		// Response is called when the clients have responded and should
		// return a string to send as a reply to the bot
		Response(payload shared.Payload) shared.Reply

		// Execute method is called when the bot command matches GetCmdName's result.
		// The bot passes the string arguments to the method.
		// The execute method must return Payload array ready to be sent to the clients.
		Execute(args ...string) ([]shared.Payload, error)

		// GetMetadata returns the metadata about the plugin.
		GetMetadata() Metadata
	}

	// Metadata is used to determine the role of a plugin registered in the PluginManager.
	Metadata struct {

		// The Type of the plugin will determine the behaviour of the server and execution of the plugin(s).
		Type string

		// The Middleware list is used to determine, if the plugin has any middleware to execute.
		// Will be skipped if the plugin itself is registered as middleware.
		Middleware []string
	}
)
