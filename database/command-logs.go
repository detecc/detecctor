package database

import (
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/shared"
	"github.com/kamva/mgm/v3"
)

func addNewCommandLog(commandLog *CommandLog) error {
	commandColl := &CommandLog{}
	return mgm.Coll(commandColl).Create(commandLog)
}

func addNewCommandResponse(commandResponse *CommandResponseLog) error {
	serviceColl := &CommandResponseLog{}
	return mgm.Coll(serviceColl).Create(commandResponse)
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

func NewCommandLog(command api.Command, payloads []shared.Payload, errs ...interface{}) (string, error) {
	commandLog := &CommandLog{
		Command:        command,
		Errors:         errs,
		PluginPayloads: payloads,
	}
	err := addNewCommandLog(commandLog)
	if err != nil {
		return "", err
	}

	return commandLog.ID.String(), nil
}
