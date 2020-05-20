package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/mocks"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var (
	testInvalidUidType = map[string]interface{}{
		"uid": strconv.Itoa(1),			// Invalid type
	}
	testUseCaseError = errors.New("usecase error")
	testMessageRequest = models.MessageRequest{
		ChatID: 1,
		Uid:    int64(security.TestUser.Uid),
		Limit:  30,
		Page:   1,
	}
	testErrorMessage = models.WorkMessage{
		Request: nil,
		Message: testUseCaseError.Error(),
		Status:  http.StatusBadRequest,
	}
	testChatRequest = models.ChatRequest{
		Uid:   int64(security.TestUser.Uid),
		Limit: 30,
		Page:  1,
	}
)

func getTestDelivery(mockUC *mocks.MockUseCase) chat.Delivery {
	return &chatDelivery{UseCase:mockUC}
}

func TestChatDelivery_GetMessages_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("PUT", "/api/chats/:id", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	vd.GetMessages(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestChatDelivery_GetMessages_InvalidCID(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("PUT", "/api/chats/:id", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	cd.GetMessages(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(-1)})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidCID, msg.Message)
}

func TestChatDelivery_GetMessages_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("PUT", "/api/chats/:id", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	cd.GetMessages(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestChatDelivery_GetMessages_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMessageRequest)
	req, err := http.NewRequest("PUT", "/api/chats/:id", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetMessagesForChat(&testMessageRequest).Return(nil, testUseCaseError)
	cd.GetMessages(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusForbidden, msg.Status)
	assert.Equal(t, testUseCaseError.Error(), msg.Message)
}

func TestChatDelivery_GetMessages_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMessageRequest)
	req, err := http.NewRequest("PUT", "/api/chats/:id", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetMessagesForChat(&testMessageRequest).Return(nil, nil)
	cd.GetMessages(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestChatDelivery_GetUsersForChat_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("GET", "/api/chats/:id/users", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	vd.GetUsersForChat(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestChatDelivery_GetUsersForChat_InvalidCID(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("GET", "/api/chats/:id/users", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	cd.GetUsersForChat(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(-1)})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidCID, msg.Message)
}

func TestChatDelivery_GetUsersForChat_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMessageRequest)
	req, err := http.NewRequest("GET", "/api/chats/:id/users", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetUsersForChat(testMessageRequest.ChatID, new(models.UserGeneralList)).Return(testErrorMessage)
	cd.GetUsersForChat(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testErrorMessage, msg)
}

func TestChatDelivery_GetUsersForChat_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMessageRequest)
	req, err := http.NewRequest("GET", "/api/chats/:id/users", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetUsersForChat(testMessageRequest.ChatID, new(models.UserGeneralList)).Return(models.WorkMessage{Message: "",})
	cd.GetUsersForChat(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestChatDelivery_GetChatList_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("POST", "/api/profile/:id/general", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	vd.GetChatList(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestChatDelivery_GetChatList_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cd := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	cd.GetChatList(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testMessageRequest.ChatID))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestChatDelivery_GetChatList_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testChatRequest)
	req, err := http.NewRequest("POST", "/api/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetUserRooms(&testChatRequest).Return(nil, testUseCaseError)
	cd.GetChatList(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testChatRequest.Uid))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testErrorMessage.Message, msg.Message)
	assert.Equal(t, http.StatusInternalServerError, msg.Status)
}

func TestChatDelivery_GetChatList_Correct1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	chatReq := testChatRequest
	chatReq.Uid = chatReq.Uid + 1
	body, _ := json.Marshal(chatReq)
	req, err := http.NewRequest("POST", "/api/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetUserRooms(&testChatRequest).Return(nil, nil)
	cd.GetChatList(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testChatRequest.Uid))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestChatDelivery_GetChatList_Correct2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	cd := getTestDelivery(mockUC)

	body, _ := json.Marshal(testChatRequest)
	req, err := http.NewRequest("POST", "/api/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetUserRooms(&testChatRequest).Return(nil, nil)
	cd.GetChatList(rr, req.WithContext(ctx), map[string]string{"id":strconv.Itoa(int(testChatRequest.Uid))})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}
