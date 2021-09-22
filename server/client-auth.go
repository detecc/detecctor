package server

import (
	"fmt"
	"github.com/detecc/detecctor/config"
	"github.com/detecc/detecctor/database"
	"github.com/detecc/detecctor/shared"
	"log"
)

// authorizeClient Try to authorize the client with the authPassword provided in the Payload.Data field
func (s *server) authorizeClient(clientId string, payload shared.Payload) error {
	log.Println("Authorizing client", clientId, "with node key", payload.ServiceNodeKey)

	if payload.Data == nil {
		return fmt.Errorf("no authentication password provided")
	}
	if payload.Data.(string) != config.GetServerConfiguration().Server.AuthPassword {
		return fmt.Errorf("invalid authentication password")
	}

	// client authorized
	err := database.AuthorizeClient(clientId, payload.ServiceNodeKey)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
