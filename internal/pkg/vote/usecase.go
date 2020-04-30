package vote

import "failless/internal/pkg/models"

type UseCase interface {
	VoteEvent(vote models.Vote) models.WorkMessage
	VoteUser(vote models.Vote) models.WorkMessage
	GetEventFollowers(eid int) (models.UserGeneralList, error)
}