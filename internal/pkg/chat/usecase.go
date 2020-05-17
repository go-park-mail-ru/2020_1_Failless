package chat

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	CreateDialogue(id1, id2 int) (int, error)
	AddNewMessage(message *forms.Message) (int, error)
	IsUserHasRoom(uid int64, cid int64) (bool, error)
	Subscribe(conn *websocket.Conn, uid int64)
	Notify(message *forms.Message)
	GetMessagesForChat(msgRequest *models.MessageRequest) (forms.MessageList, error)
	GetUserRooms(msgRequest *models.ChatRequest) (models.ChatList, error)
	GetUsersForChat(cid int64, users *models.UserGeneralList) models.WorkMessage
}

