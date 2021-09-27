package database

import (
	"github.com/detecc/detecctor/shared"
	"github.com/kamva/mgm/v3"
	"time"
)

const (
	StatusOnline       = "online"
	StatusOffline      = "offline"
	StatusUnauthorized = "unauthorized"
)

// Client contains the information of a client (Oxen Service Node).
// The object should contain the latest information of the Service Node.
type Client struct {
	mgm.DefaultModel `bson:",inline"`
	IP               string        `json:"IP" bson:"IP"`
	ClientId         string        `json:"clientId" bson:"clientId"`
	ServiceNodeKey   string        `json:"serviceNodeKey" bson:"serviceNodeKey"`
	DownTime         time.Duration `json:"downtime" bson:"downtime"`
	Status           string        `json:"status" bson:"status"`
}

// Statistics for the telegram bot.
// ActiveNodes is a number of currently active/connected Service Nodes.
// TotalNodes is a number of all known connections and is used to calculate the number of offline nodes.
// UpdateId is a number of the last known Telegram Message ID.
type Statistics struct {
	mgm.DefaultModel `bson:",inline"`
	ActiveNodes      int `json:"activeNodes" bson:"activeNodes"`
	TotalNodes       int `json:"totalNodes" bson:"totalNodes"`
	UpdateId         int `json:"updateId" bson:"updateId"`
}

// CommandLog contains information about a Command that was issued to the server, processing of the command
// and the payloads produced by the plugin as well as any errors that occurred.
type CommandLog struct {
	mgm.DefaultModel `bson:",inline"`
	Command          Command          `json:"command" bson:"command"`
	Errors           []interface{}    `json:"errors" bson:"errors"`
	PluginPayloads   []shared.Payload `json:"payloads" bson:"payloads"`
}

func (c *CommandLog) CollectionName() string {
	return "command_logs"
}

// CommandResponseLog contains the client's response to a specific command and payload.
type CommandResponseLog struct {
	mgm.DefaultModel `bson:",inline"`
	PayloadId        string        `json:"payloadId" bson:"payloadId"`
	Errors           []interface{} `json:"errors" bson:"errors"`
	PluginResponse   interface{}   `json:"pluginResponse" bson:"pluginResponse"`
}

func (c *CommandResponseLog) CollectionName() string {
	return "command_logs"
}

type Command struct {
	ChatId int64
	Name   string
	Args   []string
}

// Chat is the Telegram Chat the Bot is listening to.
type Chat struct {
	mgm.DefaultModel `bson:",inline"`
	ChatId           int64  `json:"chatId" bson:"chatId"`
	Name             string `json:"name" bson:"name"`
	IsAuthorized     bool   `json:"isAuthorized" bson:"isAuthorized"`
}

// Message is a Telegram Message that gets logged in the database.
type Message struct {
	mgm.DefaultModel `bson:",inline"`
	ChatId           int    `json:"chatId" bson:"chatId"`
	MessageId        int    `json:"messageId" bson:"messageId"`
	Content          string `json:"content" bson:"content"`
}
