package usecase

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/mocks"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	TestUId = 1
	testEventRequest = models.EventRequest{
		Uid:       1,
		Page:      0,
		Limit:     10,
		MinAmount: 3,
		MaxAmount: 15,
		Query:     "kek",
		Tags:      nil,
		Location:  models.LocationPoint{},
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
	testEventFollow = models.EventFollow{
		Uid: 1,
		Eid: 1,
	}
	repoErrorBody = errors.New("Repo error")
	repoErrorStatus = http.StatusBadRequest
	incorrectMessage = models.WorkMessage{
		Request: nil,
		Message: repoErrorBody.Error(),
		Status:  repoErrorStatus,
	}
)

func getTestUseCase(mockRep *mocks.MockRepository) event.UseCase {
	return &eventUseCase{
		rep: mockRep,
	}
}

func TestEventUseCase_GetSmallEventsByUID(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	_, err := eUC.GetSmallEventsByUID(int64(TestUId))

	assert.Equal(t, nil, err)
}

func TestEventUseCase_CreateSmallEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	YtestSmallEvent := models.SmallEvent{}
	testSmallEventForm.GetDBFormat(&YtestSmallEvent)
	testForm := testSmallEventForm
	mockRep.EXPECT().CreateSmallEvent(&YtestSmallEvent).Return(nil)
	_, err := eUC.CreateSmallEvent(&testForm)

	assert.Equal(t, nil, err)
}

func TestEventUseCase_CreateMidEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	YtestMidEvent := models.MidEvent{}
	testMidEventForm.GetDBFormat(&YtestMidEvent)
	testForm := testMidEventForm

	mockRep.EXPECT().CreateMidEvent(&YtestMidEvent).Return(repoErrorBody)
	midEvent, msg := eUC.CreateMidEvent(&testForm)

	assert.Equal(t, YtestMidEvent, midEvent)
	assert.Equal(t, incorrectMessage, msg)
}

func TestEventUseCase_CreateMidEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	YtestMidEvent := models.MidEvent{}
	testMidEventForm.GetDBFormat(&YtestMidEvent)
	testForm := testMidEventForm

	mockRep.EXPECT().CreateMidEvent(&YtestMidEvent).Return(nil)
	midEvent, msg := eUC.CreateMidEvent(&testForm)

	assert.Equal(t, YtestMidEvent, midEvent)
	assert.Equal(t, CorrectMessage, msg)
}

func TestEventUseCase_SearchEventsByUserPreferences__ZeroUid_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	events := models.MidAndBigEventList{}
	testReq := testEventRequest
	testReq.Uid = 0
	mockRep.EXPECT().GetAllMidEvents(&events.MidEvents, &testReq).Return(repoErrorStatus, repoErrorBody)

	code, err := eUC.SearchEventsByUserPreferences(&events, &testReq)

	if err == nil {
		t.Fatal()
		return
	}

	assert.Equal(t, repoErrorStatus, code)
	assert.Equal(t, repoErrorBody, err)
}

func TestEventUseCase_SearchEventsByUserPreferences_NonZeroUid_Error(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	events := models.MidAndBigEventList{}
	testReq := testEventRequest
	mockRep.EXPECT().GetMidEventsWithFollowed(&events.MidEvents, &testReq).Return(repoErrorStatus, repoErrorBody)

	code, err := eUC.SearchEventsByUserPreferences(&events, &testReq)

	if err == nil {
		t.Fatal()
		return
	}

	assert.Equal(t, repoErrorStatus, code)
	assert.Equal(t, repoErrorBody, err)
}

func TestEventUseCase_SearchEventsByUserPreferences_ZeroUid_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	events := models.MidAndBigEventList{}
	testReq := testEventRequest
	testReq.Uid = 0
	mockRep.EXPECT().GetAllMidEvents(&events.MidEvents, &testReq).Return(CorrectMessage.Status, nil)

	code, err := eUC.SearchEventsByUserPreferences(&events, &testReq)

	assert.Equal(t, CorrectMessage.Status, code)
	assert.Equal(t, nil, err)
}

func TestEventUseCase_SearchEventsByUserPreferences_NonZeroUid_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	events := models.MidAndBigEventList{}
	testReq := testEventRequest
	mockRep.EXPECT().GetMidEventsWithFollowed(&events.MidEvents, &testReq).Return(CorrectMessage.Status, nil)

	code, err := eUC.SearchEventsByUserPreferences(&events, &testReq)

	assert.Equal(t, CorrectMessage.Status, code)
	assert.Equal(t, nil, err)
}

func TestEventUseCase_UpdateSmallEvent(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().UpdateSmallEvent(&testSmallEvent).Return(CorrectMessage.Status, nil)

	code, err := eUC.UpdateSmallEvent(&testSmallEvent)

	assert.Equal(t, CorrectMessage.Status, code)
	assert.Equal(t, nil, err)
}

func TestEventUseCase_DeleteSmallEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId)).Return(repoErrorBody)

	msg := eUC.DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId))

	assert.Equal(t, incorrectMessage, msg)
}

func TestEventUseCase_DeleteSmallEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId)).Return(nil)

	msg := eUC.DeleteSmallEvent(testSmallEvent.UId, int64(testSmallEvent.EId))

	assert.Equal(t, CorrectMessage, msg)
}

func TestEventUseCase_JoinMidEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().JoinMidEvent(testEventFollow.Uid, testEventFollow.Eid).Return(repoErrorStatus, repoErrorBody)

	msg := eUC.JoinMidEvent(&testEventFollow)

	assert.Equal(t, incorrectMessage, msg)
}

func TestEventUseCase_JoinMidEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().JoinMidEvent(testEventFollow.Uid, testEventFollow.Eid).Return(CorrectMessage.Status, nil)

	msg := eUC.JoinMidEvent(&testEventFollow)

	assert.Equal(t, CorrectMessage, msg)
}

func TestEventUseCase_LeaveMidEvent_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().LeaveMidEvent(testEventFollow.Uid, testEventFollow.Eid).Return(repoErrorStatus, repoErrorBody)

	msg := eUC.LeaveMidEvent(&testEventFollow)

	assert.Equal(t, incorrectMessage, msg)
}

func TestEventUseCase_LeaveMidEvent_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	eUC := getTestUseCase(mockRep)

	mockRep.EXPECT().LeaveMidEvent(testEventFollow.Uid, testEventFollow.Eid).Return(CorrectMessage.Status, nil)

	msg := eUC.LeaveMidEvent(&testEventFollow)

	assert.Equal(t, CorrectMessage, msg)
}
