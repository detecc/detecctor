package database

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func createClient(newClient *Client) error {
	client := &Client{}
	return mgm.Coll(client).Create(newClient)
}

func updateClient(client *Client) error {
	clientCollection := &Client{}
	return mgm.Coll(clientCollection).Update(client)
}

func getClient(filter interface{}) (*Client, error) {
	client := &Client{}
	err := mgm.Coll(client).First(filter, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetClient(clientId string) (*Client, error) {
	return getClient(bson.M{"clientId": clientId})
}

func GetClientWithServiceNodeKey(serviceNodeKey string) (*Client, error) {
	return getClient(bson.M{"serviceNodeKey": serviceNodeKey})
}

// GetClients returns all clients
func GetClients() ([]Client, error) {
	var (
		serviceNode = &Client{}
		results     []Client
	)
	// find all clients
	cursor, err := mgm.Coll(serviceNode).Find(mgm.Ctx(), bson.M{})
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func IsClientAuthorized(clientId string) bool {
	sn, err := GetClient(clientId)
	if err != nil {
		return false
	}

	if sn.Status != StatusUnauthorized {
		return true
	}
	return false
}

func AuthorizeClient(clientId, serviceNodeKey string) error {
	client, err := GetClient(clientId)
	if err != nil {
		return err
	}

	client.ServiceNodeKey = serviceNodeKey

	err = updateClient(client)
	if err != nil {
		return err
	}
	return UpdateClientStatus(clientId, StatusOnline)
}

func UpdateClientStatus(clientId, status string) error {
	client, err := GetClient(clientId)
	if err != nil {
		return err
	}

	switch status {
	case StatusUnauthorized:
		client.Status = status
		break
	case StatusOffline:
		client.Status = status
		err := NodeOffline()
		if err != nil {
			return err
		}
		break
	case StatusOnline:
		client.Status = status
		err := NodeOffline()
		if err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf("client status invalid")
	}
	return updateClient(client)
}

func CreateClientIfNotExists(clientId, IP, SNKey string) *Client {
	client := &Client{
		IP:             IP,
		ClientId:       clientId,
		ServiceNodeKey: SNKey,
		DownTime:       0,
		Status:         StatusUnauthorized,
	}

	// no duplicates
	_, err := GetClient(clientId)
	if err != nil {
		log.Println("Creating client", err)
		err = createClient(client)
		if err != nil {
			return nil
		}
		return nil
	}

	return client
}
