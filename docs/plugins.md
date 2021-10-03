# Plugins

Plugins are used to customize the behaviour of the server and client to the specific needs of the user. There are a few
bundled plugins, which provide core functionality in the [Detecc-core](https://github.com/detecc/detecc-core)
repository. The plugin should have a `plugins.Handler` interface methods implemented in order to achieve desired
functionality.

All the plugins should be located in a `pluginDir`, specified in the [configuration file](../config.yaml).

The Bot will send a command, which will find the plugin and invoke the `Execute` method of the plugin. The method
returns an error and an array of `Payload`s that will be sent to different clients. After sending the message and
receiving the response from the client, the `Response` method will be invoked with the Client's
response (`Payload.Data`) and will create a `Reply` struct to send back to the Telegram Chat.

The `plugin.Handler` is shown below ([source file](../plugin/plugins.go)):

```golang
package plugin

import "github.com/detecc/detecctor/shared"

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
```

## Registering the plugins

You can register your plugins using:

```golang
package main

import "github.com/detecc/detecctor/plugin"

func init() {
	example := &YourHandlerImplementation{}
	plugin.Register("/plugin-command", example)
	//or
	GetPluginManager().AddPlugin("/plugin-command", example)
}
```

## Plugin example

```golang 
package main

import (
	"log"
	"github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
)

func init() {
	example := &Example{}
	plugin.Register("/example", example)
}

type Example struct {
	plugin.Handler
}

func (e Example) Response(payload shared.Payload) shared.Reply{
	log.Println(payload)
	return shared.Reply{
	    ChatId: payload.Id, 
	    ReplyType: shared.TypeMessage, 
	    Content: "test"
	}
}

func (e Example) Execute(args ...string) ([]shared.Payload , error) {
	log.Println(args)
	return []shared.Payload{}, nil
}
```

## Documenting plugins

Documentation of the plugins is important for further development and to better understand the logic of the plugin.

Your plugin documentation should contain:

1. What the plugin does and a brief introduction to the logic
2. The command and its arguments
3. Example(s) of the command call
4. Configuration file(s), if any are necessary
    - with a brief explanation of the attributes
    - default values, if any apply
5. The structure of the `Payload.Data`, if the plugin communicates with the client
