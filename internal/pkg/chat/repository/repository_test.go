package repository

import (
	"errors"
	"failless/internal/pkg/chat"
	"failless/internal/pkg/db/mocks"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testDBError = errors.New("db error")
	testChatReq = models.ChatRequest{
		Uid:   int64(security.TestUser.Uid),
		Limit: 10,
		Page:  0,
	}
	testChatID = int64(1)
)

func getTestRep(mockDB *mocks.MockMyDBInterface) chat.Repository {
	return &sqlChatRepository{
		db: mockDB,
	}
}

func TestSqlChatRepository_GetMessages_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	cr := getTestRep(mockDB)

	mockDB.EXPECT().Query(QuerySelectMessageHistory, testChatID, testChatReq.Uid, testChatReq.Limit, testChatReq.Page).Return(nil, testDBError)

	_, err := cr.GetRoomMessages(testChatReq.Uid, testChatID, testChatReq.Page, testChatReq.Limit)

	assert.Equal(t, testDBError, err)
}

func TestSqlChatRepository_GetUsersRooms_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	cr := getTestRep(mockDB)

	mockDB.EXPECT().Query(QuerySelectChatInfoByUserId, testChatReq.Uid).Return(nil, testDBError)

	_, err := cr.GetUsersRooms(testChatReq.Uid)

	assert.Equal(t, testDBError, err)
}

func TestSqlChatRepository_CheckRoom_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	cr := getTestRep(mockDB)

	mockDB.EXPECT().Query(QuerySelectCheckRoom, testChatID, testChatReq.Uid).Return(nil, testDBError)

	check, err := cr.CheckRoom(testChatID, testChatReq.Uid)

	assert.Equal(t, testDBError, err)
	assert.Equal(t, false, check)
}

func TestSqlChatRepository_GetUserTopMessages_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	cr := getTestRep(mockDB)

	mockDB.EXPECT().Query(QuerySelectChatsWithLastMsg, testChatReq.Uid, testChatReq.Limit, testChatReq.Page).Return(nil, testDBError)

	_, err := cr.GetUserTopMessages(testChatReq.Uid, testChatReq.Page, testChatReq.Limit)

	assert.Equal(t, testDBError, err)
}
