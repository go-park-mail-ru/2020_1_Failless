package chat

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"github.com/gorilla/websocket"
	"net/http"
)

type UseCase interface {
	CreateDialogue(id1, id2 int) (int, error)
	AddNewMessage(message *forms.Message) (int, error)
	IsUserHasRoom(uid int64, cid int64) (bool, error)
	Subscribe(conn *websocket.Conn, r *http.Request, uid int64)
	Notify(message *forms.Message)
	GetMessagesForChat(msgRequest *models.MessageRequest) ([]forms.Message, error)
	GetUserRooms(msgRequest *models.ChatRequest) ([]models.ChatMeta, error)
}

