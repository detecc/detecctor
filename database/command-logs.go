package database

import (
	"context"
	"github.com/detecc/detecctor/shared"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func addNewCommandLog(commandLog *CommandLog) error {
	commandColl := &CommandLog{}
	return mgm.Coll(commandColl).Create(commandLog)
}

func addNewCommandResponse(commandResponse *CommandResponseLog) error {
	serviceColl := &CommandResponseLog{}
	return mgm.Coll(serviceColl).Create(commandResponse)
}

func getCommandLog(filter interface{}) (*CommandLog, error) {
	serviceNode := &CommandLog{}

	err := mgm.Coll(serviceNode).First(filter, serviceNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return serviceNode, nil
}

func GetCommandLogs(filter interface{}) ([]CommandLog, error) {
	var (
		serviceNode = &CommandLog{}
		results     []CommandLog
	)

	// get service nodes with a filter
	cursor, err := mgm.Coll(serviceNode).Find(mgm.Ctx(), filter)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetAllCommandLogs() ([]CommandLog, error) {
	return GetCommandLogs(bson.M{})
}

func NewCommandResponse(payloadId string, pluginResponse interface{}, errors ...interface{}) error {
	commandResponse := &CommandResponseLog{
		PayloadId:      payloadId,
		Errors:         errors,
		PluginResponse: pluginResponse,
	}
	err := addNewCommandResponse(commandResponse)
	if err != nil {
		return err
	}
	return nil
}

func NewCommandLog(chatId int64, command string, args []string, payloads []shared.Payload, errors ...interface{}) error {
	commandLog := &CommandLog{
		Command: Command{
			Name:   command,
			Args:   args,
			ChatId: chatId,
		},
		Errors:         errors,
		PluginPayloads: payloads,
	}
	err := addNewCommandLog(commandLog)
	if err != nil {
		return err
	}
	return nil
}
