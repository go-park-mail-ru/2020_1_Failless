package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/repository"
)

type voteUseCase struct {
	rep vote.Repository
}

func GetUseCase() vote.UseCase {
	return &voteUseCase{
		rep: repository.NewSqlVoteRepository(db.ConnectToDB()),
	}
}

func (vc *voteUseCase) VoteEvent(vote models.Vote) network.Message {
	return network.Message{}
}
