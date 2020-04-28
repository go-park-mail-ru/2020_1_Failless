package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/repository"
	"log"
	"net/http"
	"strconv"

	chatRep "failless/internal/pkg/chat/repository"
)

type voteUseCase struct {
	rep vote.Repository
}

func GetUseCase() vote.UseCase {
	return &voteUseCase{
		rep: repository.NewSqlVoteRepository(db.ConnectToDB()),
	}
}

func (vc *voteUseCase) VoteUser(vote models.Vote) network.Message {
	// TODO: add check is vote already be here
	vote.Value = vc.ValidateValue(vote.Value)
	err := vc.rep.AddUserVote(vote.Uid, vote.Id, vote.Value)
	// i think that there could be an error in one case - invalid event id
	if err != nil {
		return network.Message{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	// Check for matching
	if vote.Value == 1 {
		match, _ := vc.rep.CheckMatching(vote.Uid, vote.Id)
		if match {
			log.Println("Match occured between", vote.Uid, "and", vote.Id)
			// Create dialogue
			cr := chatRep.NewSqlChatRepository(db.ConnectToDB())
			if _, err = cr.InsertDialogue(
				vote.Uid,
				vote.Id,
				2,
				"Чат#"+strconv.Itoa(vote.Id)); err != nil {
				log.Println(err)
				return network.Message{
					Request: nil,
					Message: err.Error(),
					Status:  http.StatusBadRequest,
				}
			}

			return network.Message{
				Request: nil,
				Message: "You matched with someone! Check your messages!!",
				Status:  http.StatusOK,
			}
		}
	}

	return network.Message{
		Request: nil,
		Message: "OK",
		Status:  http.StatusOK,
	}
}

func (vc *voteUseCase) VoteEvent(vote models.Vote) network.Message {
	// TODO: add check is vote already be here
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

func (vc *voteUseCase) GetEventFollowers(eid int) ([]models.UserGeneral, error) {
	return vc.rep.FindFollowers(eid)
}
