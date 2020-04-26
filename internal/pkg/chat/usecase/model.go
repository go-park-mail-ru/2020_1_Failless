package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"log"
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

func (cc *chatUseCase) IsUserHasRoom() (bool, error) {




	return false, nil
}

func (cc *chatUseCase) NotifyMembers(chatId int64) error {
	return nil
}

func (cc *chatUseCase) AddNewMessage(message *forms.Message) (int64, error) {
	// check is user has this room
	// insert message
	// notify all chat members about it
	return 0, nil
}
