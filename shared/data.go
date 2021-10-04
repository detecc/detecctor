package shared

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// constants for the Reply type
const (
	TypeMessage = 0
	TypePhoto   = 1
	TypeAudio   = 2
)

type (
	// Reply is a struct used to parse results to the ReplyChannel in Bot.
	Reply struct {
		// Each reply must contain a ChatId - a chat to reply to.
		ChatId int64
		// The ReplyType must be a constant defined in the package.
		ReplyType int
		// Content must be cast after determining the type to send to Bot.
		Content interface{}
	}

	// Payload is used to transfer data and determine the status, target client and plugin.
	Payload struct {
		// Id is used to uniquely identify the request to the response.
		Id string
		// ServiceNodeKey is used to determine the target client for the data.
		ServiceNodeKey string
		Data           interface{}
		Command        string
		Success        bool
		Error          string
	}
)

func EncodePayload(payload *Payload) (string, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	b64Payload := base64.StdEncoding.EncodeToString(data)

	return b64Payload, nil
}

func DecodePayload(data []byte, payload *Payload) error {
	jsonData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, payload)
	if err != nil {
		return err
	}
	return nil
}

// Uuid creates unique identifier.
func Uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "unknown"
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

// ParseIP separates the IP and Port of the address.
func ParseIP(addr string) (string, string) {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Println("Address not valid")
		return "", ""
	}
	return ip, port
}
