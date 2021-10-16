# Bot, Proxy and server communication

The server-proxy-bot communication occurs on channels through dedicated listeners. The message that comes from a chat
service's Bot is sent through the `chan api.ProxyMessage` channel to the bot `Proxy`. The proxy processes the message
and constructs a `Command` struct that is then sent to the server through the `chan Command` channel.

The server processes the command, executes the plugins and makes a `Reply`. The reply is sent through a `chan Reply`
channel and is passed to the proxy, which passes it to the bot.

## The Bot

The Bot interface enables the project to support multiple bots with two-way communication between the user and the
server. You can change the type of bot you want to have in the [settings](../configuration.md#configuration-file).

Check the [bot guide](adding-bots.md) for more info about the bots.

## Bot Proxy

The Bot Proxy manages the communication between the Bot and the server. It also logs and persists all the necessary data
to the database. A Bot, which is created for a specific type, is passed to the Bot Proxy. The proxy then automatically
starts the bot and its listeners.

## Command struct

The `Command` struct is created from a message that is sent to the bot and must have a `/` prefix. The message is split
by spaces; the first element in the split array is `Command.Name`, while the other elements of the array represent
the `Command.Args` (arguments).

The arguments are passed to the plugin which handles the `Command`. The `ChatId` is used to map the plugin response to
the chat.

```go
package bot

type Command struct {
	Name   string
	Args   []string
	ChatId string
}
```

### Building the command(s)

To easily construct the commands from the bot, use the `CommandBuilder`.

```go
package example

import "github.com/detecc/detecctor/bot/api"

func CreateCommand() {
	builder := api.NewMessageBuilder()
	replyMessage := builder.WithId("chatId1").FromUser("userName").WithMessage("messageId", "messageText").Build()
}
```

## Reply struct

The `Reply` struct is used in server-to-bot communication. The `ChatId` is used to reply to a specific chat.
The `ReplyType` represents one of the constants and specifies the type of `Content`. The `Content`
is a generic representation of the data, usually produced by the plugin (e.g. an Image, Text, Audio file, etc.)

```go
package bot

// constants for the Reply type
const (
	TypeMessage = 0
	TypePhoto   = 1
	TypeAudio   = 2
)

type Reply struct {
	ChatId    string
	ReplyType string
	Content   interface{}
}
```

### Building the replies

To easily construct the replies for the bot, use the `ReplyBuilder`.

```go
package example

import "github.com/detecc/detecctor/bot/api"

func CreateReply() {
	builder := api.NewReplyBuilder()
	replyMessage1 := builder.TypeMessage().ForChat("chatId1").WithContent("sampleContent").Build()
	replyMessage2 := builder.TypePhoto().ForChat("chatId2").WithContent(nil).Build()
}
```