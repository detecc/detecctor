package shared

import (
	"encoding/base64"
	"encoding/json"
)

// SetError indicate that something went wrong with the payload
func (payload *Payload) SetError(err error) {
	if err != nil {
		payload.Success = false
		payload.Error = err.Error()
	}
}

// NewPayload creates a new payload with the options provided
func NewPayload(opts ...PayloadOption) Payload {
	var payload = &Payload{
		Id:             "",
		ServiceNodeKey: "",
		Data:           nil,
		Command:        "",
		Success:        false,
		Error:          "",
	}

	for _, opt := range opts {
		opt(payload)
	}

	return *payload
}

// ForClient add a target client for the payload
func ForClient(serviceNodeKey string) PayloadOption {
	return func(payload *Payload) {
		payload.ServiceNodeKey = serviceNodeKey
	}
}

// ForCommand add a command which the payload should execute or is originating from
func ForCommand(command string) PayloadOption {
	return func(payload *Payload) {
		payload.Command = command
	}
}

// WithData sets the payload data
func WithData(data interface{}) PayloadOption {
	return func(payload *Payload) {
		payload.Data = data
	}
}

// Successful set that the payload is successful
func Successful() PayloadOption {
	return func(payload *Payload) {
		payload.Error = ""
		payload.Success = true
	}
}

// WithError sets an error for the payload
func WithError(err error) PayloadOption {
	return func(payload *Payload) {
		payload.Error = err.Error()
		payload.Success = false
	}
}

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
