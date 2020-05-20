package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/mocks"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	testEventRequest = models.EventRequest{
		Uid:       1,
		Page:      1,
		Limit:     10,
		UserLimit: 15,
		Query:     "kek",
		Tags:      nil,
		Location:  models.LocationPoint{},
		MinAge:    18,
		MaxAge:    100,
		Men:       true,
		Women:     true,
	}
	testSmallEvent = models.SmallEvent{
		EId:    1,
		UId:    1,
		Title:  "title",
		Descr:  "about",
		TagsId: nil,
		Date:   time.Time{},
		Photos: nil,
	}
	testSmallEventForm = forms.SmallEventForm{
		Uid:    testSmallEvent.EId,
		Title:  testSmallEvent.Title,
		Descr:  testSmallEvent.Descr,
		TagsId: testSmallEvent.TagsId,
		Date:   testSmallEvent.Date,
		Photos: nil,
	}
	testEventFollow = models.EventFollow{
		Uid: 1,
		Eid: 1,
	}
	testInvalidUidType = map[string]interface{}{
		"uid": strconv.Itoa(1),			// Invalid type
	}
)

func getTestDelivery(mockUC *mocks.MockUseCase) event.Delivery {
	return &eventDelivery{UseCase:mockUC}
}

func TestEventDelivery_GetEventsByKeyWords_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	mockVoteBody := map[string]interface{}{
		"page": strconv.Itoa(testEventRequest.Page),			// Invalid type
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.GetEventsByKeyWords(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_GetEventsByKeyWords_IncorrectPage(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	mockVoteBody := map[string]interface{}{
		"page": 0,			// Invalid type
		"query": testEventRequest.Query,
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().InitEventsByKeyWords(new(models.EventList), testEventRequest.Query, 1)
	ed.GetEventsByKeyWords(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_GetEventsByKeyWords_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().InitEventsByKeyWords(new(models.EventList), testEventRequest.Query, testEventRequest.Page).Return(http.StatusInternalServerError, errors.New("error"))
	ed.GetEventsByKeyWords(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "error", msg.Message)
}

func TestEventDelivery_GetEventsByKeyWords_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().InitEventsByKeyWords(new(models.EventList), testEventRequest.Query, testEventRequest.Page)
	ed.GetEventsByKeyWords(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_GetSearchEvents_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	mockVoteBody := map[string]interface{}{
		"page": strconv.Itoa(testEventRequest.Page),			// Invalid type
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.GetSearchEvents(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_GetSearchEvents_IncorrectPage(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	searchReq := testEventRequest
	searchReq.Page = 0
	body, _ := json.Marshal(searchReq)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &testEventRequest)
	ed.GetSearchEvents(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_GetSearchEvents_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &testEventRequest).Return(http.StatusInternalServerError, errors.New("error"))
	ed.GetSearchEvents(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "error", msg.Message)
}

func TestEventDelivery_GetSearchEvents_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &testEventRequest)
	ed.GetSearchEvents(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_GetSmallEvents_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("GET", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.GetSmallEvents(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_GetSmallEvents_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("GET", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetSmallEventsByUID(int64(1)).Return(nil, errors.New("error in usecase"))
	ed.GetSmallEvents(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, "error in usecase", msg.Message)
}

func TestEventDelivery_GetSmallEvents_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("GET", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().GetSmallEventsByUID(int64(1)).Return(nil, nil)
	ed.GetSmallEvents(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_CreateSmallEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.CreateSmallEvent(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_CreateSmallEvent_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_CreateSmallEvent_InvalidForm(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	eventReq := testSmallEventForm
	eventReq.Uid = 0
	body, _ := json.Marshal(eventReq)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, forms.MessageSmallEventValidationFailed, msg.Message)
}

func TestEventDelivery_CreateSmallEvent_InvalidPhotos(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	eventReq := testSmallEventForm
	eventReq.Photos = []forms.EImage{{ImgBase64:""}, {ImgBase64:""}}
	body, _ := json.Marshal(eventReq)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, images.MessageImageValidationFailed, msg.Message)
}

func TestEventDelivery_CreateSmallEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEventForm)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().CreateSmallEvent(&testSmallEventForm).Return(models.SmallEvent{}, errors.New("error in usecase"))
	ed.CreateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, "error in usecase", msg.Message)
}

func TestEventDelivery_CreateSmallEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEventForm)
	req, err := http.NewRequest("POST", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().CreateSmallEvent(&testSmallEventForm).Return(models.SmallEvent{}, nil)
	ed.CreateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_UpdateSmallEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventRequest)
	req, err := http.NewRequest("PUT", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.UpdateSmallEvent(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_UpdateSmallEvent_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	mockVoteBody := map[string]interface{}{
		"uid": strconv.Itoa(1),			// Invalid type
	}
	body, _ := json.Marshal(mockVoteBody)
	req, err := http.NewRequest("PUT", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.UpdateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_UpdateSmallEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("PUT", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().UpdateSmallEvent(&testSmallEvent).Return(http.StatusInternalServerError, errors.New("error in usecase"))
	ed.UpdateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "error in usecase", msg.Message)
}

func TestEventDelivery_UpdateSmallEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("GET", "/api/srv/events/small", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().UpdateSmallEvent(&testSmallEvent).Return(http.StatusOK, nil)
	ed.UpdateSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_JoinMiddleEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("POST", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	ed.JoinMiddleEvent(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_JoinMiddleEvent_IncorrectEidInUrl(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("POST", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = "kek"
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.JoinMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorRetrievingEidFromUrl, msg.Message)
}

func TestEventDelivery_JoinMiddleEvent_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.JoinMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_JoinMiddleEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("POST", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().JoinMidEvent(&testEventFollow).Return(models.WorkMessage{
		Request: nil,
		Message: "error in usecase",
		Status:  http.StatusInternalServerError,
	})
	ed.JoinMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "error in usecase", msg.Message)
}

func TestEventDelivery_JoinMiddleEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("POST", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().JoinMidEvent(&testEventFollow)
	ed.JoinMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}
