package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"log"
	"net/http"
)

type chatUseCase struct {
	Rep chat.Repository
}

func GetUseCase() chat.UseCase {
	return &chatUseCase{
		Rep: repository.NewSqlChatRepository(db.ConnectToDB()),
	}
}

func (cc *chatUseCase) CreateDialogue(id1, id2 int) (int, error) {
	chatId, err := cc.CreateDialogue(id1, id2)

	if err != nil {
		log.Println(err)
		return -1, nil
	}

	return chatId, nil
}

func (cc *chatUseCase) IsUserHasRoom(uid int64, cid int64) (bool, error) {
	return cc.Rep.CheckRoom(cid, uid)
}

func (cc *chatUseCase) NotifyMembers(chatId int64) error {
	// TODO: implement it
	return nil
}

func (cc *chatUseCase) AddNewMessage(message *forms.Message) (int, error) {
	// check is user has this room
	has, err := cc.IsUserHasRoom(message.Uid, message.ChatID)
	if err != nil {
		log.Println("AddNewMessage: error - ", err.Error())
		return http.StatusInternalServerError, err
	}
	if !has {
		log.Println("AddNewMessage: client - room not found: ", message.Uid, message.ChatID)
		return http.StatusNotFound, nil
	}
	// insert message
	msgID, err := cc.Rep.AddMessageToChat(message, nil)
	if err != nil {
		log.Println("AddNewMessage: error while AddMessageToChat -  ", err.Error())
		return http.StatusInternalServerError, err
	}

	// TODO: check it
	message.ULocalID = msgID
	// notify all chat members about it
	err = cc.NotifyMembers(message.ChatID)
	if err != nil {
		// TODO: create cool handler or remove it
		return http.StatusServiceUnavailable, nil
	}
	log.Println("AddNewMessage: OK")
	return http.StatusOK, nil
}

func (cc *chatUseCase) GetMessagesForChat(msgRequest *models.MessageRequest) ([]forms.Message, error) {
	has, err := cc.IsUserHasRoom(msgRequest.Uid, 0)
	if err != nil || !has {
		return nil, err
	}
	return cc.Rep.GetRoomMessages(msgRequest.Uid, msgRequest.ChatID, msgRequest.Page, msgRequest.Limit)
}

func (cc *chatUseCase) GetUserRooms(msgRequest *models.ChatRequest) ([]models.ChatMeta, error) {
	return cc.Rep.GetUserTopMessages(msgRequest.Uid, msgRequest.Page, msgRequest.Limit)
}
