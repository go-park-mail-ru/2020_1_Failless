package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/repository"
	"failless/internal/pkg/db"
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
