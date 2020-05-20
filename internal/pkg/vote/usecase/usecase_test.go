package usecase

import (
	"errors"
	chatMocks "failless/internal/pkg/chat/mocks"
	"failless/internal/pkg/models"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var (
	testRepoError = errors.New("repo error")
	testVote = models.Vote{
		Uid:   1,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}
)

func getTestUseCase(mockRep *mocks.MockRepository, chatMockRep *chatMocks.MockRepository) vote.UseCase {
	return &voteUseCase{
		Rep: mockRep,
		chatRep: chatMockRep,
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
	vc := getTestUseCase(mockRep, chatMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value)
	mockRep.EXPECT().CheckMatching(testVote.Uid, testVote.Id)

	vc.VoteUser(testVote)
}

func TestVoteUseCase_VoteUser_InvalidUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep, chatMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(errors.New("DB error"))

	msg := vc.VoteUser(testVote)
	assert.Equal(t, msg.Status, http.StatusBadRequest)
}

func TestVoteUseCase_VoteUser_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep, chatMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(testRepoError)

	msg := vc.VoteUser(testVote)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepoError.Error(), msg.Message)
}

func TestVoteUseCase_VoteUser_Correct1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep, chatMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(nil)
	mockRep.EXPECT().CheckMatching(testVote.Uid, testVote.Id).Return(false, nil)

	msg := vc.VoteUser(testVote)

	assert.Equal(t, CorrectMessage, msg)
}

func TestVoteUseCase_VoteUser_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	chatMockRep := chatMocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep, chatMockRep)

	mockRep.EXPECT().AddUserVote(testVote.Uid, testVote.Id, testVote.Value).Return(nil)
	mockRep.EXPECT().CheckMatching(testVote.Uid, testVote.Id).Return(true, nil)
	chatMockRep.EXPECT().InsertDialogue(testVote.Uid, testVote.Id, 2, "Чат#"+strconv.Itoa(testVote.Id)).Return(int64(0), testRepoError)

	msg := vc.VoteUser(testVote)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepoError.Error(), msg.Message)
}
