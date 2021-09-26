package database

import (
	"github.com/kamva/mgm/v3"
	"time"
)

const (
	StatusOnline       = "online"
	StatusOffline      = "offline"
	StatusUnauthorized = "unauthorized"
)

// Client object contains some basic information of a client.
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

// ServiceNode contains the information of an Oxen Service Node.
// This entry is used for storing historical data about the Node.
// ClientId is a key connected to the Client.
// ServiceNodeKey is a Service Node key.
// ServiceNodeStatus is a response from the Client.
type ServiceNode struct {
	mgm.DefaultModel  `bson:",inline"`
	ClientId          string `json:"clientId" bson:"clientId"`
	ServiceNodeKey    string `json:"serviceNodeKey" bson:"serviceNodeKey"`
	ServiceNodeStatus string `json:"status" bson:"status"`
}

// Chat is the Telegram Chat the Bot is listening to.
type Chat struct {
	mgm.DefaultModel `bson:",inline"`
	ChatId           int64          `json:"chatId" bson:"chatId"`
	Name             string         `json:"name" bson:"name"`
	IsAuthorized     bool           `json:"isAuthorized" bson:"isAuthorized"`
	Subscriptions    []Subscription `json:"subscriptions" bson:"subscriptions"`
}

// Subscription is a filter used for subscribing to a client messages.
// If the chat/user is subscribed to all nodes and all topics, there should only be one entry with both subNode and subCommand values equal to "*".
// Else, there are separate entries with values, "*" meaning all.
// Example entry: subNode: "*", subCommand:"/ping" -> meaning subscribe to the /ping command on all nodes.
type Subscription struct {
	Node    string `json:"nodeId" bson:"nodeId"`
	Command string `json:"command" bson:"command"`
}

// Message is a Telegram Message that gets logged in the database.
type Message struct {
	mgm.DefaultModel `bson:",inline"`
	ChatId           int    `json:"chatId" bson:"chatId"`
	MessageId        int    `json:"messageId" bson:"messageId"`
	Content          string `json:"content" bson:"content"`
}
