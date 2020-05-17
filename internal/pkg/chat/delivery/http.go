package delivery

import (
	"failless/internal/pkg/chat/usecase"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

//func SendMessage(w http.ResponseWriter, r *http.Request, ps map[string]string) {
//	uid := security.CheckCredentials(w, r)
//	if uid < 0 {
//		return
//	}
//	cid := int64(0)
//	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
//		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
//		return
//	}
//
//	decoder := json.NewDecoder(r.Body)
//	var message forms.Message
//	err := decoder.Decode(&message)
//	if err != nil {
//		network.Jsonify(w, "Error within parse json", http.StatusBadRequest)
//		return
//	}
//	if int64(uid) != message.Uid {
//		message.Uid = int64(uid)
//		log.Println("SendMessage: warn - uid from token is not equal message.Uid")
//	}
//	message.ChatID = cid
//	uc := usecase.GetUseCase()
//
//	if code, err := uc.AddNewMessage(&message); err != nil {
//		network.GenErrorCode(w, r, err.Error(), code)
//		return
//	}
//	network.Jsonify(w, models.WorkMessage{
//		Request: nil,
//		Message: "OK",
//		Status:  http.StatusOK,
//	}, http.StatusOK)
//}

func GetMessages(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("GetMessages: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
		return
	}

	var request models.MessageRequest
	err := json.UnmarshalFromReader(r.Body, &request)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	request.ChatID = cid
	uc := usecase.GetUseCase()
	messages, err := uc.GetMessagesForChat(&request)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusForbidden)
		return
	}
	log.Println("OK")
	network.Jsonify(w, messages, http.StatusOK)
}

func GetUsersForChat(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("GetMessages: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}
	cid := int64(0)
	if cid = network.GetIdFromRequest(w, r, ps); cid < 0 {
		network.GenErrorCode(w, r, "url cid is incorrect", http.StatusBadRequest)
		return
	}

	var users models.UserGeneralList
	uc := usecase.GetUseCase()
	message := uc.GetUsersForChat(cid, &users)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, users, message.Status)
}

func GetChatList(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Print("GetChatList: ")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var request models.ChatRequest
	err := json.UnmarshalFromReader(r.Body, &request)
	if err != nil {
		log.Println("error within parse json - ", err.Error())
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	if int64(uid) != request.Uid {
		request.Uid = int64(uid)
		log.Println("warn - uid from token is not equal request.Uid")
	}

	uc := usecase.GetUseCase()
	chatList, err := uc.GetUserRooms(&request)
	if err != nil {
		log.Println("error while GetUserRooms - ", err.Error())
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("OK")
	network.Jsonify(w, chatList, http.StatusOK)
}
