package usecase

import (
	"errors"
	eventMocks "failless/internal/pkg/event/mocks"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	testRepError = errors.New("error in repo")
	testUserRequest = models.UserRequest{
		Uid:      security.TestUser.Uid,
		Page:     1,
		Limit:    10,
		Query:    "kek",
		Tags:     nil,
		Location: models.LocationPoint{},
		MinAge:   18,
		MaxAge:   100,
		Men:      true,
		Women:    true,
	}
	testUser = models.UserGeneral{
		Uid:      security.TestUser.Uid,
		Name:     security.TestUser.Name,
		Photos:   nil,
		About:    "about",
		Birthday: time.Time{},
		Gender:   0,
		TagsId:   nil,
	}
	testUserWithPass = models.User{
		Uid:      security.TestUser.Uid,
		Name:     security.TestUser.Name,
		Phone:    security.TestUser.Phone,
		Email:    security.TestUser.Email,
		Password: nil,
	}
	testSignForm = forms.SignForm{
		Uid:      security.TestUser.Uid,
		Name:     security.TestUser.Name,
		Phone:    security.TestUser.Phone,
		Email:    security.TestUser.Email,
		Password: "qwerty1234",
	}
	testGeneralForm = forms.GeneralForm{
		SignForm: testSignForm,
		Tags:     nil,
		Avatar:   forms.EImage{},
		Photos:   nil,
		Gender:   0,
		About:    "",
		Rating:   0,
		Location: models.LocationPoint{},
		Birthday: time.Time{},
	}
)

func getTestUseCase(mockRep *mocks.MockRepository, mockEventRep *eventMocks.MockRepository) user.UseCase {
	return &UserUseCase{
		Rep: mockRep,
		eventRep: mockEventRep,
	}
}

func TestUserUseCase_InitUsersByUserPreferences_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetValidTags().Return(nil, testRepError)

	code, err := uUC.InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)

	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_InitUsersByUserPreferences_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetValidTags().Return(nil, nil)
	mockRep.EXPECT().GetRandomFeedUsers(testUserRequest.Uid, testUserRequest.Limit, testUserRequest.Page).Return(nil, testRepError)

	code, err := uUC.InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_InitUsersByUserPreferences_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetValidTags().Return(nil, nil)
	mockRep.EXPECT().GetRandomFeedUsers(testUserRequest.Uid, testUserRequest.Limit, testUserRequest.Page).Return(nil, nil)

	code, err := uUC.InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, nil, err)
}

func TestUserUseCase_GetUserSubscriptions_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	subs := models.MidAndBigEventList{}
	eventMockRep.EXPECT().GetSubscriptionMidEvents(&subs.MidEvents, testUserRequest.Uid).Return(http.StatusBadRequest, testRepError)

	msg := uUC.GetUserSubscriptions(&subs, testUserRequest.Uid)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetUserSubscriptions_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	subs := models.MidAndBigEventList{}
	eventMockRep.EXPECT().GetSubscriptionMidEvents(&subs.MidEvents, testUserRequest.Uid).Return(http.StatusOK, nil)

	msg := uUC.GetUserSubscriptions(&subs, testUserRequest.Uid)

	assert.Equal(t, CorrectMessage.Status, msg.Status)
	assert.Equal(t, CorrectMessage.Message, msg.Message)
}

func TestUserUseCase_GetFeedResultsFor_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	feed := models.FeedResults{models.FeedPost{}}
	eventMockRep.EXPECT().GetSmallEventsForUser(&feed[0].OnwEvents.SmallEvents, testUserRequest.Uid).Return(http.StatusBadRequest, testRepError)

	var users []models.UserGeneral
	users = append(users, models.UserGeneral{Uid: testUserRequest.Uid})
	_, msg := uUC.GetFeedResultsFor(testUserRequest.Uid, &users)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetFeedResultsFor_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	feed := models.FeedResults{models.FeedPost{}}
	eventMockRep.EXPECT().GetSmallEventsForUser(&feed[0].OnwEvents.SmallEvents, testUserRequest.Uid)
	eventMockRep.EXPECT().GetOwnMidEventsWithAnotherUserFollowed(&feed[0].OnwEvents.MidEvents, testUserRequest.Uid, testUserRequest.Uid).Return(http.StatusBadRequest, testRepError)

	var users []models.UserGeneral
	users = append(users, models.UserGeneral{Uid: testUserRequest.Uid})
	_, msg := uUC.GetFeedResultsFor(testUserRequest.Uid, &users)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetFeedResultsFor_Incorrect3(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	feed := models.FeedResults{models.FeedPost{}}
	eventMockRep.EXPECT().GetSmallEventsForUser(&feed[0].OnwEvents.SmallEvents, testUserRequest.Uid)
	eventMockRep.EXPECT().GetOwnMidEventsWithAnotherUserFollowed(&feed[0].OnwEvents.MidEvents, testUserRequest.Uid, testUserRequest.Uid)
	eventMockRep.EXPECT().GetSubscriptionMidEventsWithAnotherUserFollowed(&feed[0].Subscriptions.MidEvents, testUserRequest.Uid, testUserRequest.Uid).Return(http.StatusBadRequest, testRepError)

	var users []models.UserGeneral
	users = append(users, models.UserGeneral{Uid: testUserRequest.Uid})
	_, msg := uUC.GetFeedResultsFor(testUserRequest.Uid, &users)

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetFeedResultsFor_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	feed := models.FeedResults{models.FeedPost{}}
	eventMockRep.EXPECT().GetSmallEventsForUser(&feed[0].OnwEvents.SmallEvents, testUserRequest.Uid)
	eventMockRep.EXPECT().GetOwnMidEventsWithAnotherUserFollowed(&feed[0].OnwEvents.MidEvents, testUserRequest.Uid, testUserRequest.Uid)
	eventMockRep.EXPECT().GetSubscriptionMidEventsWithAnotherUserFollowed(&feed[0].Subscriptions.MidEvents, testUserRequest.Uid, testUserRequest.Uid)

	var users []models.UserGeneral
	users = append(users, models.UserGeneral{Uid: testUserRequest.Uid})
	_, msg := uUC.GetFeedResultsFor(testUserRequest.Uid, &users)

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_UpdateUserAbout_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().UpdateUserAbout(testUser.Uid, testUser.About).Return(CorrectMessage)

	msg := uUC.UpdateUserAbout(testUser.Uid, testUser.About)

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_UpdateUserTags_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	testTags := []int{1,2}
	mockRep.EXPECT().UpdateUserTags(testUser.Uid, testTags).Return(CorrectMessage)

	msg := uUC.UpdateUserTags(testUser.Uid, testTags)

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_GetUserInfo_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetProfileInfo(testGeneralForm.Uid).Return(models.JsonInfo{}, testRepError)

	code, err := uUC.GetUserInfo(&testGeneralForm)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_GetUserInfo_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetProfileInfo(testGeneralForm.Uid)
	mockRep.EXPECT().GetUserByUID(testGeneralForm.Uid).Return(models.User{}, testRepError)

	code, err := uUC.GetUserInfo(&testGeneralForm)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_GetUserInfo_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetProfileInfo(testGeneralForm.Uid)
	mockRep.EXPECT().GetUserByUID(testGeneralForm.Uid)

	code, err := uUC.GetUserInfo(&testGeneralForm)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, nil, err)
}

func TestUserUseCase_UpdateUserPhotos_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().UpdateUserPhotos(testGeneralForm.Uid, new([]string)).Return(CorrectMessage)

	msg := uUC.UpdateUserPhotos(testGeneralForm.Uid, new(forms.EImageList))

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_UpdateUserBase_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUserByUID(testGeneralForm.Uid).Return(models.User{}, testRepError)

	code, err := uUC.UpdateUserBase(&testSignForm)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_GetSmallEventsForUser_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	eventMockRep.EXPECT().GetSmallEventsForUser(new(models.SmallEventList), testUser.Uid).Return(http.StatusInternalServerError, testRepError)

	msg := uUC.GetSmallEventsForUser(new(models.SmallEventList), testUser.Uid)

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetSmallEventsForUser_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	eventMockRep.EXPECT().GetSmallEventsForUser(new(models.SmallEventList), testUser.Uid)

	msg := uUC.GetSmallEventsForUser(new(models.SmallEventList), testUser.Uid)

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_GetUserOwnEvents_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	midEventList := models.MidEventList{}
	eventMockRep.EXPECT().GetOwnMidEvents(&midEventList, testUser.Uid).Return(http.StatusInternalServerError, testRepError)

	msg := uUC.GetUserOwnEvents(new(models.OwnEventsList), testUser.Uid)

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetUserOwnEvents_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	midEventList := models.MidEventList{}
	eventMockRep.EXPECT().GetOwnMidEvents(&midEventList, testUser.Uid)
	smallEventList := models.SmallEventList{}
	eventMockRep.EXPECT().GetSmallEventsForUser(&smallEventList, testUser.Uid).Return(http.StatusInternalServerError, testRepError)

	msg := uUC.GetUserOwnEvents(new(models.OwnEventsList), testUser.Uid)

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, testRepError.Error(), msg.Message)
}

func TestUserUseCase_GetUserOwnEvents_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	eventMockRep := eventMocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mocks.NewMockRepository(mockCtrl), eventMockRep)

	midEventList := models.MidEventList{}
	eventMockRep.EXPECT().GetOwnMidEvents(&midEventList, testUser.Uid)
	smallEventList := models.SmallEventList{}
	eventMockRep.EXPECT().GetSmallEventsForUser(&smallEventList, testUser.Uid)

	msg := uUC.GetUserOwnEvents(new(models.OwnEventsList), testUser.Uid)

	assert.Equal(t, CorrectMessage, msg)
}

func TestUserUseCase_FillFormIfExist_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUserByPhoneOrEmail(testUserWithPass.Phone, testUserWithPass.Email).Return(models.User{Uid:-1}, nil)

	code, err := uUC.FillFormIfExist(&testUserWithPass)

	if err == nil {
		t.Fatal()
		return
	}

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, MessageUserDoesntExist, err.Error())
}

func TestUserUseCase_FillFormIfExist_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUserByPhoneOrEmail(testUserWithPass.Phone, testUserWithPass.Email).Return(models.User{Uid:1}, testRepError)

	code, err := uUC.FillFormIfExist(&testUserWithPass)

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, testRepError, err)
}

func TestUserUseCase_FillFormIfExist_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	uUC := getTestUseCase(mockRep, eventMocks.NewMockRepository(mockCtrl))

	mockRep.EXPECT().GetUserByPhoneOrEmail(testUserWithPass.Phone, testUserWithPass.Email).Return(models.User{Uid:1}, nil)

	code, err := uUC.FillFormIfExist(&testUserWithPass)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, nil, err)
}
