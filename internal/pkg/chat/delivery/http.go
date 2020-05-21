package delivery

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

type chatDelivery struct {
	UseCase chat.UseCase
}

func GetDelivery() chat.Delivery {
	return &chatDelivery{
		UseCase:usecase.GetUseCase(),
	}
}

func (cd *chatDelivery) GetMessages(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("GetMessages: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, network.MessageInvalidCID, http.StatusBadRequest)
		return
	}

	var request models.MessageRequest
	err := json.UnmarshalFromReader(r.Body, &request)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	request.ChatID = cid
	messages, err := cd.UseCase.GetMessagesForChat(&request)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusForbidden)
		return
	}
	log.Println("OK")
	network.Jsonify(w, messages, http.StatusOK)
}

func (cd *chatDelivery) GetUsersForChat(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("GetMessages: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, network.MessageInvalidCID, http.StatusBadRequest)
		return
	}

	var users models.UserGeneralList
	message := cd.UseCase.GetUsersForChat(cid, &users)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, users, message.Status)
}

func (cd *chatDelivery) GetChatList(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	log.Print("GetChatList: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var request models.ChatRequest
	err := json.UnmarshalFromReader(r.Body, &request)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	if int64(uid) != request.Uid {
		request.Uid = int64(uid)
		log.Println("warn - uid from token is not equal request.Uid")
	}

	chatList, err := cd.UseCase.GetUserRooms(&request)
	if err != nil {
		log.Println("error while GetUserRooms - ", err.Error())
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("OK")
	network.Jsonify(w, chatList, http.StatusOK)
}
