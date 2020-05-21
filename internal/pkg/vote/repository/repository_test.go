package repository

import (
	"errors"
	"failless/internal/pkg/db/mocks"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	testDBError = errors.New("db error")
	testUserVote = models.Vote{
		Uid:   security.TestUser.Uid,
		Id:    security.TestUser.Uid + 1,
		Value: 1,
		Date:  time.Time{},
	}
)

func getTestRep(mockDB *mocks.MockMyDBInterface) vote.Repository {
	return &sqlVoteRepository{
		db: mockDB,
	}
}

func TestSqlVoteRepository_AddUserVote_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	vr := getTestRep(mockDB)

	mockDB.EXPECT().Exec(QueryInsertUserVote, testUserVote.Uid, testUserVote.Id, testUserVote.Value).Return(pgx.CommandTag(""), testDBError)
	err := vr.AddUserVote(testUserVote.Uid, testUserVote.Id, testUserVote.Value)

	assert.Equal(t, testDBError, err)
}

func TestSqlVoteRepository_AddUserVote_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	vr := getTestRep(mockDB)

	mockDB.EXPECT().Exec(QueryInsertUserVote, testUserVote.Uid, testUserVote.Id, testUserVote.Value).Return(pgx.CommandTag(""), nil)
	err := vr.AddUserVote(testUserVote.Uid, testUserVote.Id, testUserVote.Value)

	assert.Equal(t, nil, err)
}
