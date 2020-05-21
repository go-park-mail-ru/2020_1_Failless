package repository

import (
	"errors"
	"failless/internal/pkg/db/mocks"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	testUser = models.User{
		Uid:      security.TestUser.Uid,
		Name:     security.TestUser.Name,
		Phone:    security.TestUser.Phone,
		Email:    security.TestUser.Email,
		Password: nil,
	}
	testLocation = models.LocationPoint{}
	testDBError = errors.New("db error")
	testErrorMessage = models.WorkMessage{
		Request: nil,
		Message: testDBError.Error(),
		Status:  http.StatusInternalServerError,
	}
	testFeedRequest = models.UserRequest{
		Uid:      security.TestUser.Uid,
		Page:     1,
		Limit:    10,
		Query:    "",
		Tags:     nil,
		Location: models.LocationPoint{},
		MinAge:   0,
		MaxAge:   0,
		Men:      false,
		Women:    false,
	}
)

func getTestRep(mockDB *mocks.MockMyDBInterface) user.Repository {
	return &sqlUserRepository{
		db: mockDB,
	}
}

func TestSqlUserRepository_SetUserLocation(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	mockDB.EXPECT().Exec(QueryUpdateUserLocation, testLocation.Latitude, testLocation.Longitude, testUser.Uid)
	err := ur.SetUserLocation(testUser.Uid, testLocation)

	assert.NoError(t, err)
}

func TestSqlUserRepository_UpdateUserRating(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	tmpRating := float32(1)
	mockDB.EXPECT().Exec(QueryUpdateUserRating, tmpRating, testUser.Uid)
	err := ur.UpdateUserRating(testUser.Uid, tmpRating)

	assert.NoError(t, err)
}

func TestSqlUserRepository_UpdateUserTags(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	tmpTags := []int{1,2}

	// Incorrect
	mockDB.EXPECT().Exec(QueryUpdateUserTags, tmpTags, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), testDBError)
	msg := ur.UpdateUserTags(testUser.Uid, tmpTags)
	assert.Equal(t, testErrorMessage, msg)

	// Correct
	mockDB.EXPECT().Exec(QueryUpdateUserTags, tmpTags, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), nil)
	msg = ur.UpdateUserTags(testUser.Uid, tmpTags)
	assert.Equal(t, CorrectMessage, msg)
}

func TestSqlUserRepository_UpdateUserAbout(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	tmpAbout := "about"

	// Incorrect
	mockDB.EXPECT().Exec(QueryUpdateUserAbout, tmpAbout, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), testDBError)
	msg := ur.UpdateUserAbout(testUser.Uid, tmpAbout)
	assert.Equal(t, testErrorMessage, msg)

	// Correct
	mockDB.EXPECT().Exec(QueryUpdateUserAbout, tmpAbout, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), nil)
	msg = ur.UpdateUserAbout(testUser.Uid, tmpAbout)
	assert.Equal(t, CorrectMessage, msg)
}

func TestSqlUserRepository_DeleteUser(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	tmpEmail := "mail@mail.ru"
	mockDB.EXPECT().Exec(QueryDeleteUserByEmail, tmpEmail)
	err := ur.DeleteUser(tmpEmail)

	assert.NoError(t, err)
}

func TestSqlUserRepository_UpdateUserPhotos(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	tmpPhotos := new([]string)
	*tmpPhotos = append(*tmpPhotos, "photo1")

	// Incorrect
	mockDB.EXPECT().Exec(QueryUpdateUserPhotos, &tmpPhotos, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), testDBError)
	msg := ur.UpdateUserPhotos(testUser.Uid, tmpPhotos)
	assert.Equal(t, testErrorMessage, msg)

	// Correct
	mockDB.EXPECT().Exec(QueryUpdateUserPhotos, &tmpPhotos, testUser.Uid).Return(pgx.CommandTag("UPDATE 0 1"), nil)
	msg = ur.UpdateUserPhotos(testUser.Uid, tmpPhotos)
	assert.Equal(t, CorrectMessage, msg)
}

func TestSqlUserRepository_GetValidTags(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	// Incorrect
	mockDB.EXPECT().Query(QuerySelectTags).Return(nil, testDBError)
	_, err := ur.GetValidTags()
	assert.Equal(t, testDBError, err)
}

func TestSqlUserRepository_GetRandomFeedUsers(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	// Incorrect args 1
	users, err := ur.GetRandomFeedUsers(testFeedRequest.Uid, testFeedRequest.Limit, 0)
	assert.Nil(t, users)
	assert.Equal(t, errors.New("Page number can't be less than 1\n"), err)

	// Incorrect args 2
	users, err = ur.GetRandomFeedUsers(testFeedRequest.Uid, 0, testFeedRequest.Page)
	assert.Nil(t, users)
	assert.Equal(t, errors.New("Page number can't be less than 1\n"), err)

	// Incorrect
	mockDB.EXPECT().Query(QueryWithVotedUsersIncomplete+QuerySelectUserInfoIncomplete+QueryConditionFeedIncomplete, testFeedRequest.Uid, testFeedRequest.Limit).Return(nil, testDBError)
	_, err = ur.GetRandomFeedUsers(testFeedRequest.Uid, testFeedRequest.Limit, testFeedRequest.Page)
	assert.Equal(t, testDBError, err)
}

func TestSqlUserRepository_GetUsersForChat(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	ur := getTestRep(mockDB)

	chatID := int64(1)

	// Incorrect
	mockDB.EXPECT().Query(QueryWithChatMembersIncomplete+QuerySelectUserInfoIncomplete+QueryConditionChatMembersIncomplete, chatID).Return(nil, testDBError)
	msg := ur.GetUsersForChat(chatID, new(models.UserGeneralList))
	assert.Equal(t, testErrorMessage, msg)
}
