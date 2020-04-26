package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
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

func (cc *chatUseCase) IsUserHasRoom(int64) (bool, error) {

	return false, nil
}

func (cc *chatUseCase) NotifyMembers(chatId int64) error {
	// TODO: implement it
	return nil
}

func (cc *chatUseCase) AddNewMessage(message *forms.Message) (int, error) {
	// check is user has this room
	has, err := cc.IsUserHasRoom(message.Uid)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !has {
		return http.StatusNotFound, nil
	}
	// insert message
	msgID, err := cc.Rep.AddMessageToChat(message, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: check it
	message.ULocalID = msgID
	// notify all chat members about it
	err = cc.NotifyMembers(message.ChatId)
	if err != nil {
		// TODO: create cool handler or remove it
		return http.StatusServiceUnavailable, nil
	}

	return http.StatusOK, nil
}
