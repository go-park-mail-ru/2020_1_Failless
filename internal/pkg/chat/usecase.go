package chat

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	CreateDialogue(id1, id2 int) (int, error)
	AddNewMessage(message *forms.Message) (int, error)
	IsUserHasRoom(uid int64, cid int64) (bool, error)
	NotifyMembers(chatId int64) error
	GetMessagesForChat(msgRequest *models.MessageRequest) ([]forms.Message, error)
	GetUserRooms(msgRequest *models.ChatRequest) ([]models.ChatMeta, error)
}
