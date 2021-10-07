package server

import (
	"container/list"
	"github.com/Allenxuxu/gev"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/shared"
	"sync"
)

const (
	AuthCommand          = "/auth"
	SubscribeCommand     = "/sub"
	UnSubscribeCommand   = "/unsub"
	LanguageCommand      = "/language"
	LanguageCommandShort = "/lang"
)

// Server TCP/Websockets server.
// The server handles the incoming messages from the Bot and executes the command.
// The command execution is handled by plugins that are loaded in the cache.
// When a command is received, it forwards the command to a specific client. After the command was executed in the client,
// the server sends the command result back to the Bot.
type server struct {
	conn         *list.List
	mu           sync.RWMutex
	server       *gev.Server
	botChannel   chan bot.Command
	replyChannel chan shared.Reply
}
