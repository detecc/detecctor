package shared

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
)

type (
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

	PayloadOption func(*Payload)
)

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
