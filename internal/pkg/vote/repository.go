package vote

import "failless/internal/pkg/models"

type Repository interface {
	AddEventVote(uid int, eid int, value int8) error
	AddUserVote(uid int, id int, value int8) error
	FindFollowers(eid int) ([]models.UserGeneral, error)
	CheckMatching(uid, id int) (bool, error)
}
