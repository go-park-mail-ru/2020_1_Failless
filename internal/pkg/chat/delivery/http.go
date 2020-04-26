package delivery

import (
	"encoding/json"
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/network"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var message forms.Message
	err := decoder.Decode(&message)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
		return
	}

	message.ChatId = cid
	uc := usecase.GetUseCase()
	if code, err := uc.AddNewMessage(&message); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var message forms.Message
	err := decoder.Decode(&message)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
		return
	}

	message.ChatId = cid
	uc := usecase.GetUseCase()
	if code, err := uc.UpdateUserBase(&message); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
}

func GetChatList(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var message forms.Message
	err := decoder.Decode(&message)
	if err != nil {
		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
		return
	}

	message.ChatId = cid
	uc := usecase.GetUseCase()
	if code, err := uc.UpdateUserBase(&message); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
}