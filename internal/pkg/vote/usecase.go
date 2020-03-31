package vote

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
)

type UseCase interface {
	VoteEvent(vote models.Vote) network.Message
}