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
	TestEventRequest = models.EventRequest{
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
	testMidEventForm = forms.MidEventForm{
		AdminId: 1,
		Title:   "title",
		Descr:   "about",
		TagsId:  nil,
		Date:    time.Time{},
		Photos:  nil,
		Limit:   10,
		Public:  false,
	}
	testMessageUseCaseOk = models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  0,
	}
	testMessageUseCaseError = models.WorkMessage{
		Request: nil,
		Message: "error in usecase",
		Status:  http.StatusInternalServerError,
	}
)

func getTestDelivery(mockUC *mocks.MockUseCase) event.Delivery {
	return &eventDelivery{UseCase:mockUC}
}

func TestEventDelivery_GetSearchEvents_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	mockVoteBody := map[string]interface{}{
		"page": strconv.Itoa(TestEventRequest.Page),			// Invalid type
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

	searchReq := TestEventRequest
	searchReq.Page = 0
	body, _ := json.Marshal(searchReq)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &TestEventRequest)
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

	body, _ := json.Marshal(TestEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &TestEventRequest).Return(http.StatusInternalServerError, errors.New("error"))
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

	body, _ := json.Marshal(TestEventRequest)
	req, err := http.NewRequest("POST", "/api/srv/events/search", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	mockUC.EXPECT().SearchEventsByUserPreferences(new(models.MidAndBigEventList), &TestEventRequest)
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

	body, _ := json.Marshal(TestEventRequest)
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

	body, _ := json.Marshal(TestEventRequest)
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

	body, _ := json.Marshal(TestEventRequest)
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
	assert.Equal(t, forms.MessageEventValidationFailed, msg.Message)
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

	body, _ := json.Marshal(TestEventRequest)
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

func TestEventDelivery_DeleteSmallEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("DELETE", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.DeleteSmallEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_DeleteSmallEvent_IncorrectEidInUrl(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("DELETE", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = "kek"
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.DeleteSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorRetrievingEidFromUrl, msg.Message)
}

func TestEventDelivery_DeleteSmallEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("DELETE", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testSmallEvent.EId)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId)).Return(models.WorkMessage{
		Request: nil,
		Message: "error in usecase",
		Status:  http.StatusInternalServerError,
	})
	ed.DeleteSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "error in usecase", msg.Message)
}

func TestEventDelivery_DeleteSmallEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSmallEvent)
	req, err := http.NewRequest("DELETE", "/api/srv/events/small/:eid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testSmallEvent.EId)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId))
	ed.DeleteSmallEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestEventDelivery_CreateMiddleEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMidEventForm)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.CreateMiddleEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_CreateMiddleEvent_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateMiddleEvent(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_CreateMiddleEvent_InvalidForm(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	eventReq := testMidEventForm
	eventReq.AdminId = 0
	body, _ := json.Marshal(eventReq)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateMiddleEvent(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, forms.MessageEventValidationFailed, msg.Message)
}

func TestEventDelivery_CreateMiddleEvent_InvalidPhotos(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	eventReq := testMidEventForm
	eventReq.Photos = []forms.EImage{{ImgBase64:""}, {ImgBase64:""}}
	body, _ := json.Marshal(eventReq)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.CreateMiddleEvent(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, images.MessageImageValidationFailed, msg.Message)
}

func TestEventDelivery_CreateMiddleEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMidEventForm)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().CreateMidEvent(&testMidEventForm).Return(models.MidEvent{}, testMessageUseCaseError)
	ed.CreateMiddleEvent(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestEventDelivery_CreateMiddleEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testMidEventForm)
	req, err := http.NewRequest("POST", "/api/srv/events/mid", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().CreateMidEvent(&testMidEventForm).Return(models.MidEvent{}, testMessageUseCaseOk)
	ed.CreateMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseOk.Status, msg.Status)
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

func TestEventDelivery_LeaveMiddleEvent_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.LeaveMiddleEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestEventDelivery_LeaveMiddleEvent_IncorrectEidInUrl(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = "kek"
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.LeaveMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorRetrievingEidFromUrl, msg.Message)
}

func TestEventDelivery_LeaveMiddleEvent_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.LeaveMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestEventDelivery_LeaveMiddleEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().LeaveMidEvent(&testEventFollow).Return(testMessageUseCaseError)
	ed.LeaveMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestEventDelivery_LeaveMiddleEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testEventFollow)
	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid/member", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ps := make(map[string]string)
	ps["eid"] = strconv.Itoa(testEventFollow.Eid)
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().LeaveMidEvent(&testEventFollow)
	ed.LeaveMiddleEvent(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseOk.Status, msg.Status)
}


// Not implemented

func TestEventDelivery_GetMiddleEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/events/mid/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.GetMiddleEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_UpdateMiddleEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/srv/events/mid/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.UpdateMiddleEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_DeleteMiddleEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("DELETE", "/api/srv/events/mid/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.DeleteMiddleEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_CreateBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("POST", "/api/srv/events/big", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.CreateBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_GetBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/events/big/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.GetBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_UpdateBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/srv/events/big/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.UpdateBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_DeleteBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("DELETE", "/api/srv/events/big/:eid", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.DeleteBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_AddVisitorForBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("POST", "/api/srv/events/big/:eid/visitor", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.AddVisitorForBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}

func TestEventDelivery_RemoveVisitorForBigEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	req, err := http.NewRequest("DELETE", "/api/srv/events/big/:eid/visitor", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.RemoveVisitorForBigEvent(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotImplemented, msg.Status)
	assert.Equal(t, http.StatusText(http.StatusNotImplemented), msg.Message)
}
