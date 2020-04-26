package chat

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type Repository interface {
	InsertDialogue(id1, id2 int) (int, error)

	GetUsersRooms(uid int64) ([]models.ChatRoom, error)
	CheckRoom(cid int64, uid int64) (bool, error)
	AddMessageToChat(msg *forms.Message, relatedChats []int64) (int64, error)
	GetUserTopMessages(uid int64, page, limit int) ([]forms.Message, error)
	GetRoomMessages(uid, cid int64, page, limit int) ([]forms.Message, error)
}
