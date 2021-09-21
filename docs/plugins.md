# Plugins

## About

Plugins are used to customize the behaviour of the server and client the specific needs of the user. There are a few
bundled plugins, which provide core functionality and can be replaced or customized entirely. The plugin should have
a `plugins.Handler` interface methods implemented in order to achieve desired functionality.

All the plugins should be located at a `pluginDir`, specified in the [configuration file](../Readme.md#configuration).

The Bot will send a command, which will find the plugin and invoke the `Execute` method of the plugin. The method
returns an error and an array of `Payload`s that will be sent to different clients. After sending the message and
receiving the response from the client, the `Response` method will be invoked with the Client's
response (`Payload.Data`) and will create a `Reply` struct to send back to the Telegram Chat.

The `plugins.Handler` is shown below:

```golang
import "github.com/detecc/detecctor/shared"

type Handler interface {
// Response is called when the clients have responded and should
// return a string to send as a reply to the bot
Response(payload shared.Payload) shared.Reply
// Execute method is called when the bot command matches GetCmdName's result.
// The bot passes the string arguments to the method.
//The execute method must return Payload array ready to be sent to the clients.
Execute(args ...string) ([]shared.Payload, error)
}
```

## Example

Example plugin file:

```golang 
package main

import (
	"log"
	"github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
	"github.com/detecc/detecctor/bot"
)

func init() {
	example := &Example{}
	plugin.Register("/example", example)
}

type Example struct {
	plugin.Handler
}

func (e Example) Response(payload shared.Payload) bot.Reply{
	log.Println(payload)
	return bot.Reply{
	    ChatId: payload.Id, 
	    ReplyType: bot.TYPE_MESSAGE, 
	    Content: "test"
	}
}

func (e Example) Execute(args ...string) ([]shared.Payload , error) {
	log.Println(args)
	return []shared.Payload{}, nil
}
```