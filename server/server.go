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
	plugin2 "github.com/detecc/detecctor/server/plugin"
	"github.com/detecc/detecctor/shared"
	"log"
	"sync"
	"time"
)

var srv *server
var once = sync.Once{}

// Start a new TCP/WS server.
func Start(botChannel chan bot.Command, replyChannel chan shared.Reply) error {
	var err error
	serverConfig := config.GetServerConfiguration()

	if botChannel == nil {
		return fmt.Errorf("bot channel is nil")
	}
	if replyChannel == nil {
		return fmt.Errorf("reply channel is nil")
	}

	once.Do(func() {
		srv = &server{
			conn:         list.New(),
			mu:           sync.RWMutex{},
			botChannel:   botChannel,
			replyChannel: replyChannel,
		}

		address := fmt.Sprintf("%s:%d", serverConfig.Server.Host, serverConfig.Server.Port)
		srv.server, err = gev.NewServer(srv, gev.Address(address))
		if err != nil {
			log.Fatal(err)
		}
	})

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
		log.Println("Received command:", command)
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

		notificationMessage := fmt.Sprintf("Client %s went offline at %s", client.ServiceNodeKey, time.Now().Format(time.RFC1123))
		log.Println(notificationMessage)

		data := make(map[string]interface{})
		data["ServiceNodeKey"] = client.ServiceNodeKey
		data["Time"] = time.Now().Format(time.RFC1123)
		message := MakeTranslationMap("ClientDisconnected", nil, data)

		//notify the user(s) the node went down
		for _, chat := range chats {
			s.replyToChat(chat.ChatId, message, shared.TypeMessage)
		}
	}
}

// validateCommand check if the chat is authorized to perform a command.
func (s *server) validateCommand(command bot.Command) error {
	if !s.isChatAuthorized(command.ChatId) && command.Name != "/auth" {
		return fmt.Errorf("chat is not authorized")
	}

	return nil
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
