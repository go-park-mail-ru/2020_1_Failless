package chat

import "failless/internal/pkg/forms"

type UseCase interface {
	CreateDialogue(id1, id2 int) (int, error)
	AddNewMessage(message *forms.Message) (int64, error)
	IsUserHasRoom() (bool, error)
	NotifyMembers(chatId int64) error
}
