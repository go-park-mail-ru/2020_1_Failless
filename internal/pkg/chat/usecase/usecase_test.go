package usecase

import (
	"errors"
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/mocks"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	userMocks "failless/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	testRepoError = errors.New("repo error")
	testChatRequest = models.ChatRequest{
		Uid:   int64(security.TestUser.Uid),
		Limit: 30,
		Page:  1,
	}
	testMessage = forms.Message{
		Mid:      1,
		Uid:      int64(security.TestUser.Uid),
		ULocalID: 1,
		IsShown:  false,
		ChatID:   1,
		Text:     "",
		Date:     time.Time{},
	}
	testMsgReq = models.MessageRequest{
		ChatID: testMessage.ChatID,
		Uid:    testMessage.Uid,
		Limit:  30,
		Page:   1,
	}
	UID = testMessage.Uid
	CID = testMessage.ChatID
	testWorkMessage = models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  0,
	}
)

func getTestUseCase(mockRep *mocks.MockRepository, userMockRep *userMocks.MockRepository) chat.UseCase {
	return &chatUseCase{
		Rep: mockRep,
		userRep: userMockRep,
	}
}

func TestChatUseCase_CreateDialogue_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().InsertDialogue(1, 2, 2, "").Return(int64(1), testRepoError)

	cID, err := cc.CreateDialogue(1, 2)

	assert.Equal(t, -1, cID)
	assert.Equal(t, testRepoError, err)
}

func TestChatUseCase_CreateDialogue_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().InsertDialogue(1, 2, 2, "").Return(int64(1), nil)

	cID, err := cc.CreateDialogue(1, 2)

	assert.Equal(t, 1, cID)
	assert.Equal(t, nil, err)
}

func TestChatUseCase_IsUserHasRoom(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(UID, CID).Return(true, nil)

	answer, err := cc.IsUserHasRoom(UID, CID)

	assert.Equal(t, true, answer)
	assert.Equal(t, nil, err)
}

func TestChatUseCase_Subscribe_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUsersRooms(UID).Return(nil, testRepoError)

	cc.Subscribe(new(websocket.Conn), UID)
}

func TestChatUseCase_AddNewMessage_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(testMessage.ChatID, testMessage.Uid).Return(false, testRepoError)

	code, err := cc.AddNewMessage(&testMessage)

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, testRepoError, err)
}

func TestChatUseCase_AddNewMessage_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(testMessage.ChatID, testMessage.Uid).Return(false, nil)

	code, err := cc.AddNewMessage(&testMessage)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, nil, err)
}

func TestChatUseCase_AddNewMessage_Incorrect3(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(testMessage.ChatID, testMessage.Uid).Return(true, nil)
	mockRep.EXPECT().AddMessageToChat(&testMessage, nil).Return(int64(1), testRepoError)

	code, err := cc.AddNewMessage(&testMessage)

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, testRepoError, err)
}

func TestChatUseCase_AddNewMessage_Incorrect4(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(testMessage.ChatID, testMessage.Uid).Return(true, nil)
	mockRep.EXPECT().AddMessageToChat(&testMessage, nil).Return(int64(1), nil)

	code, err := cc.AddNewMessage(&testMessage)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, nil, err)
}

func TestChatUseCase_GetMessagesForChat_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(UID, CID).Return(true, testRepoError)

	msgList, err := cc.GetMessagesForChat(&testMsgReq)

	assert.Equal(t, testRepoError, err)
	assert.Equal(t, forms.MessageList(nil), msgList)
}

func TestChatUseCase_GetMessagesForChat(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().CheckRoom(UID, CID).Return(true, nil)
	mockRep.EXPECT().GetRoomMessages(UID, CID, testMsgReq.Page, testMsgReq.Limit).Return(nil, nil)

	msgList, err := cc.GetMessagesForChat(&testMsgReq)

	assert.Equal(t, nil, err)
	assert.Equal(t, forms.MessageList(nil), msgList)
}

func TestChatUseCase_GetUserRooms(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep, userMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUserTopMessages(UID, testMsgReq.Page, testMsgReq.Limit).Return(nil, nil)

	msgList, err := cc.GetUserRooms(&testChatRequest)

	assert.Equal(t, nil, err)
	assert.Equal(t, models.ChatList(nil), msgList)
}

func TestChatUseCase_GetUsersForChat(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	usrMockRep := userMocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mocks.NewMockRepository(mockCtrl), usrMockRep)

	usrMockRep.EXPECT().GetUsersForChat(CID, new(models.UserGeneralList)).Return(testWorkMessage)

	msg := cc.GetUsersForChat(CID, new(models.UserGeneralList))

	assert.Equal(t, testWorkMessage, msg)
}
