package usecase

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"net/http"
	"testing"
	"time"
)

func getTestUseCase(mockRep *mocks.MockRepository) vote.UseCase {
	return &voteUseCase{
		rep: mockRep,
	}
}

func TestGetUseCase(t *testing.T) {

}

func TestVoteUseCase_VoteEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep)

	var testVote = models.Vote{
		Uid:   0,
		Id:    0,
		Value: -1,
		Date:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockRep.EXPECT().AddEventVote(testVote.Uid, testVote.Id, testVote.Value).Return(nil)

	res := vc.VoteEvent(testVote)
	assert.Equal(t,
		network.Message{
			Request: nil,
			Message: "OK",
			Status:  http.StatusOK,
		},
		res)
}

func TestVoteUseCase_ValidateValue(t *testing.T) {
	vc := new(voteUseCase)
	assert.Equal(t, int8(1), vc.ValidateValue(1))
	assert.Equal(t, int8(1), vc.ValidateValue(math.MaxInt8))
	assert.Equal(t, int8(-1), vc.ValidateValue(0))
	assert.Equal(t, int8(-1), vc.ValidateValue(-1))
	assert.Equal(t, int8(-1), vc.ValidateValue(-math.MaxInt8))
}

func TestVoteUseCase_GetEventFollowers(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep)

	mockRep.EXPECT().FindFollowers(0)

	vc.GetEventFollowers(0)
}
