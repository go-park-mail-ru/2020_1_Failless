package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/repository"
	"net/http"
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
	vote.Value = vc.ValidateValue(vote.Value)
	err := vc.rep.AddEventVote(vote.Uid, vote.Id, vote.Value)
	// i think that there could be an error in one case - invalid event id
	if err != nil {
		return network.Message{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return network.Message{
		Request: nil,
		Message: "OK",
		Status:  http.StatusOK,
	}
}

func (vc *voteUseCase) ValidateValue(value int8) int8 {
	switch {
	case value >= 1:
		return 1
	case value <= 1:
		return -1
	}
	return 0
}
