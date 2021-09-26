package server

import (
	"container/list"
	"fmt"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/detecc/detecctor/bot"
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/config"
	"github.com/detecc/detecctor/database"
	plugin2 "github.com/detecc/detecctor/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"sync"
	"time"
)

// Start a new TCP/WS server.
func Start(botChannel chan bot.Command, replyChannel chan shared.Reply) error {
	serverConfig := config.GetServerConfiguration()
	var err error
	if botChannel == nil {
		return fmt.Errorf("bot channel is nil")
	}
	if replyChannel == nil {
		return fmt.Errorf("reply channel is nil")
	}

	srv := &server{
		conn:         list.New(),
		mu:           sync.RWMutex{},
		botChannel:   botChannel,
		replyChannel: replyChannel,
	}

	address := fmt.Sprintf("%s:%d", serverConfig.Server.Host, serverConfig.Server.Port)
	srv.server, err = gev.NewServer(srv, gev.Address(address))
	if err != nil {
		return err
	}

	plugin2.GetPluginManager().LoadPlugins()

	srv.start()
	return nil
}

// start server
func (s *server) start() {
	go s.listenForCommands()
	s.server.Start()
}

// stop server
func (s *server) stop() {
	s.server.Stop()
}

// ListenForCommands listen for incoming bot commands
func (s *server) listenForCommands() {
	for command := range s.botChannel {
		log.Println("received command:", command)
		go s.handleCommand(command)
	}
}

// OnConnect handle the client connection
func (s *server) OnConnect(c *connection.Connection) {
	log.Println("Client connected:", c.PeerAddr())
	s.storeClient(c)
}

// OnMessage handle the incoming client messages
func (s *server) OnMessage(c *connection.Connection, ctx interface{}, data []byte) interface{} {
	var payloadData shared.Payload
	log.Println("Received message from client:", string(data))

	err := shared.DecodePayload(data, &payloadData)
	if err != nil {
		log.Println(err)
		c.ShutdownWrite()
		return nil
	}

	s.handleMessage(c, payloadData)
	return nil
}

// OnClose Handle the closing of the connection
func (s *server) OnClose(c *connection.Connection) {
	log.Println("Closing connection from", c.PeerAddr())
	e := c.Context().(*list.Element)

	clientId, _ := shared.ParseIP(c.PeerAddr())
	client, err := database.GetClient(clientId)
	// update the client's status
	database.UpdateClientStatus(clientId, database.StatusOffline)

	// remove the connection from the connection list
	s.mu.Lock()
	s.conn.Remove(e)
	cache.Memory().Delete(c.PeerAddr())
	s.mu.Unlock()

	if err == nil {
		chats, err := database.GetChats()
		if err != nil {
			log.Println(err)
			return
		}
		//notify the user(s) the node went down
		notificationMessage := fmt.Sprintf("Client %s went offline at %s", client.ServiceNodeKey, time.Now().Format(time.RFC1123))
		log.Println(notificationMessage)
		for _, chat := range chats {
			s.replyToChat(chat.ChatId, notificationMessage, shared.TypeMessage)
		}
	}
}

// validateCommand check if the chat is authorized to perform a command.
func (s *server) validateCommand(command bot.Command) error {
	if !s.isChatAuthorized(command.ChatId) && command.Name != "/auth" {
		return fmt.Errorf("chat is not authorized")
	}

	// todo include validation middleware here

	return nil
}

// handleCommand handles the invocation of the Plugin.Execute method and sends the payloads produced to the designated clients.
func (s *server) handleCommand(command bot.Command) {
	cmdErr := s.validateCommand(command)
	if cmdErr != nil {
		s.replyToChat(command.ChatId, "You are not authorized to send this command", shared.TypeMessage)
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

//executePlugin executes the plugin associated with the command and sends a message to the client(s).
func (s *server) executePlugin(command bot.Command) {
	//check if the plugin is authorized
	plugin, err := plugin2.GetPluginManager().GetPlugin(command.Name)
	if err != nil {
		log.Println("Plugin with command", command.Name, "doesnt exist")
		s.replyToChat(command.ChatId, fmt.Sprintf("%s unsupported command", command.Name), shared.TypeMessage)
		return
	}

	// invoke the Plugin.Execute method
	payloads, err := plugin.Execute(command.Args...)
	if err != nil {
		log.Println("plugin produced an error:", err)
		s.replyToChat(command.ChatId, "command could not be executed.", shared.TypeMessage)
		return
	}

	if payloads != nil {
		// send the payloads to the clients
		for i, payload := range payloads {
			generatePayloadId(&payload, command.ChatId)
			log.Println(i, payload)
			messageErr := s.sendMessage(payload)
			if messageErr != nil {
				couldNotSendMessage := fmt.Sprintf("could  not send message to %s: %v", payload.ServiceNodeKey, messageErr)
				log.Println(couldNotSendMessage)
				s.replyToChat(command.ChatId, couldNotSendMessage, shared.TypeMessage)
			}
		}
	}
}

func (s *server) replyToChat(chatId int64, content interface{}, contentType int) {
	if contentType < 0 {
		contentType = shared.TypeMessage
	}

	s.replyChannel <- shared.Reply{
		ChatId:    chatId,
		ReplyType: contentType,
		Content:   content,
	}
}

// getConnection returns a connection pointer stored in memory based on clientId
func (s *server) getConnection(serviceNodeKey string) (*connection.Connection, error) {
	client, err := database.GetClientWithServiceNodeKey(serviceNodeKey)
	if err != nil {
		return nil, err
	}
	conn, ok := cache.Memory().Get(client.ClientId)
	if !ok {
		return nil, fmt.Errorf("Could not find a connected client with Service Node Key: %s ", serviceNodeKey)
	}
	return conn.(*connection.Connection), nil
}
