package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func createServiceNodeEntry(serviceNodeEntry *ServiceNode) error {
	serviceColl := &ServiceNode{}
	return mgm.Coll(serviceColl).Create(serviceNodeEntry)
}

func getServiceNode(filter interface{}) (*ServiceNode, error) {
	serviceNode := &ServiceNode{}

	err := mgm.Coll(serviceNode).First(filter, serviceNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return serviceNode, nil
}

func FindServiceNodeWithClientId(clientId string) (*ServiceNode, error) {
	return getServiceNode(bson.M{"clientId": clientId})
}

func FindServiceNodeWithKey(serviceNodeKey string) (*ServiceNode, error) {
	return getServiceNode(bson.M{"serviceNodeKey": serviceNodeKey})
}

func getServiceNodes(filter interface{}) ([]ServiceNode, error) {
	var (
		serviceNode = &ServiceNode{}
		results     []ServiceNode
	)

	// get service nodes with a filter
	cursor, err := mgm.Coll(serviceNode).Find(mgm.Ctx(), filter)
	if err = cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetServiceNodes() ([]ServiceNode, error) {
	return getServiceNodes(bson.M{})
}

func GetServiceNodeEntriesWithKey(serviceNodeKey string) ([]ServiceNode, error) {
	return getServiceNodes(bson.M{"serviceNodeKey": serviceNodeKey})
}

func NewServiceNode(clientId, serviceNodeKey, status string) *ServiceNode {
	serviceNode := &ServiceNode{
		ClientId:          clientId,
		ServiceNodeKey:    serviceNodeKey,
		ServiceNodeStatus: status,
	}
	err := createServiceNodeEntry(serviceNode)
	if err != nil {
		return nil
	}
	return serviceNode
}

func UpdateServiceNode(clientId, serviceNodeKey, status string) error {
	serviceNode, _ := FindServiceNodeWithClientId(clientId)
	updatedSN := serviceNode
	updatedSN.ServiceNodeKey = serviceNodeKey
	updatedSN.ServiceNodeStatus = status
	return mgm.Coll(serviceNode).Update(updatedSN)
}
