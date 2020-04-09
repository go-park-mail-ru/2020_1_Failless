package network

import (
	"bytes"
	"encoding/json"
)

func DecodeToMsg(body *bytes.Buffer) (Message, error) {
	var msg Message
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}


