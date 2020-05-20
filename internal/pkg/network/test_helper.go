package network

import (
	"bytes"
	"encoding/json"
	"failless/internal/pkg/models"
)

const (
	MessageErrorParseJSON = "Error within parse json"
	MessageErrorAuthRequired = "auth required"
	MessageErrorRetrievingEidFromUrl = "Error in retrieving eid from url"
	MessageInvalidID = "Incorrect id"
	MessageValidationFailed = "Validation failed"
	MessageSuccessfulLogout = "Successful logout"
)

func DecodeToMsg(body *bytes.Buffer) (models.WorkMessage, error) {
	var msg models.WorkMessage
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&msg)
	if err != nil {
		return models.WorkMessage{}, err
	}
	return msg, nil
}
