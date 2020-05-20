package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote"
	"failless/internal/pkg/vote/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func getTestDelivery(mockUC *mocks.MockUseCase) vote.Delivery {
	return &voteDelivery{UseCase:mockUC}
}

func TestVoteUser(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	vd := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/users/:vote", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	vd.VoteUser(rr, req, ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestVoteDelivery_VoteUser_IncorrectUidInBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	vd := getTestDelivery(mockUC)

	testVote := models.Vote{
		Uid:   1,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}
	mockVoteBody := map[string]interface{}{
		"uid": testVote.Uid + 1,			// wrong Uid
		"id": testVote.Id,
		"value": testVote.Value,
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("GET", "/api/srv/users/:vote", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mail@mail.ru",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)
	vd.VoteUser(rr, req.WithContext(ctx), ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, msg.Status, http.StatusBadRequest)
	assert.Equal(t, msg.Message, MessageInvalidUidInBody)
}

func TestVoteDelivery_VoteUser_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	vd := getTestDelivery(mockUC)

	testVote := models.Vote{
		Uid:   1,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}
	mockVoteBody := map[string]interface{}{
		"uid": testVote.Uid,
		"id": testVote.Id,
		"value": testVote.Value,
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("GET", "/api/srv/users/:vote", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	user := security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mail@mail.ru",
		Name:  "mrTester",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, user)
	correctMessage := models.WorkMessage{
		Request: nil,
		Message: "OK",
		Status:  http.StatusOK,
	}
	mockUC.EXPECT().VoteUser(testVote).Return(correctMessage)
	vd.VoteUser(rr, req.WithContext(ctx), ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, msg.Status, correctMessage.Status)
	assert.Equal(t, msg.Message, correctMessage.Message)
}

func TestVoteDelivery_VoteUser_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	vd := getTestDelivery(mockUC)

	testVote := models.Vote{
		Uid:   1,
		Id:    2,
		Value: 1,
		Date:  time.Time{},
	}
	mockVoteBody := map[string]interface{}{
		"uid": testVote.Uid,
		"id": strconv.Itoa(testVote.Id),		// Invalid type
		"value": testVote.Value,
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("GET", "/api/srv/users/:vote", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	rr := httptest.NewRecorder()

	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mail@mail.ru",
		Name:  "mrTester",
	})
	vd.VoteUser(rr, req.WithContext(ctx), ps)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, msg.Status, http.StatusBadRequest)
	assert.Equal(t, msg.Message, network.MessageErrorParseJSON)
}
