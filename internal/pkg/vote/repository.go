package vote

import "failless/internal/pkg/models"

type Repository interface {
	AddEventVote(uid int, eid int, value int8) error
	FindFollowers(eid int) ([]models.UserGeneral, error)
}
