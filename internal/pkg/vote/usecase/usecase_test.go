package usecase

import (
	"errors"
	"failless/internal/pkg/models"
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

func TestVoteUseCase_ValidateValue(t *testing.T) {
	vc := new(voteUseCase)
	assert.Equal(t, int8(1), vc.ValidateValue(1))
	assert.Equal(t, int8(1), vc.ValidateValue(math.MaxInt8))
	assert.Equal(t, int8(-1), vc.ValidateValue(0))
	assert.Equal(t, int8(-1), vc.ValidateValue(-1))
	assert.Equal(t, int8(-1), vc.ValidateValue(-math.MaxInt8))
}

func TestVoteUseCase_VoteUser(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep)

	testVote := models.Vote{
		Uid:   1,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value)
	mockRep.EXPECT().CheckMatching(testVote.Uid, testVote.Id)

	vc.VoteUser(testVote)
}

func TestVoteUseCase_VoteUser_InvalidUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep)

	testVote := models.Vote{
		Uid:   0,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(errors.New("DB error"))

	msg := vc.VoteUser(testVote)
	assert.Equal(t, msg.Status, http.StatusBadRequest)
}

//func TestVoteUseCase_VoteUser_Check(t *testing.T) {
//	// Create mock
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	mockRep := mocks.NewMockRepository(mockCtrl)
//	vc := getTestUseCase(mockRep)
//
//	testVote := models.Vote{
//		Uid:   1,
//		Id:    2,
//		Value: 1,
//		Date:  time.Time{},
//	}
//
//	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(nil)
//	mockRep.EXPECT().CheckMatching(testVote.Uid, testVote.Id).Return(true, nil)
//
//	msg := vc.VoteUser(testVote)
//	assert.Equal(t, msg.Status, http.StatusOK)
//	assert.Equal(t, msg.Message, MessageMatchHappened)
//}
