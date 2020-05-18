package vote

import (
	"failless/internal/pkg/models"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	VoteEvent(vote models.Vote) models.WorkMessage
	VoteUser(vote models.Vote) models.WorkMessage
	GetEventFollowers(eid int) (models.UserGeneralList, error)
	Subscribe(conn *websocket.Conn, uid int64)
}