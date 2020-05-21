package vote

import (
	"failless/internal/pkg/models"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	VoteUser(vote models.Vote) models.WorkMessage
	Subscribe(conn *websocket.Conn, uid int64)
}