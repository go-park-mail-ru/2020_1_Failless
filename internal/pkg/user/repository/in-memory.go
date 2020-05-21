package repository

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
	"net/http"
)

type userRepository struct {
}


func NewUserRepository() user.Repository {
	return &userRepository{}
}

func (ur *userRepository) GetUserByUID(uid int) (models.User, error) {
	return models.User{}, nil
}

func (ur *userRepository) GetUserByPhoneOrEmail(phone string, email string) (models.User, error) {
	return models.User{}, nil
}

// Private method
//func (ur *userRepository) getUser(sqlStatement string, args ...interface{}) (user models.User, err error) {
//	return models.User{}, nil
//}

func (ur *userRepository) AddNewUser(user *models.User) error {
	return nil
}

func (ur *userRepository) SetUserLocation(uid int, point models.LocationPoint) error {
	return nil
}

func (ur *userRepository) UpdateUserRating(uid int, rating float32) error {
	return nil
}

func (ur *userRepository) UpdateUserTags(uid int, tagIDs []int) models.WorkMessage {
	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusNotImplemented),
		Status:  http.StatusNotImplemented,
	}
}

func (ur *userRepository) UpdateUserAbout(uid int, about string) models.WorkMessage {
	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusNotImplemented),
		Status:  http.StatusNotImplemented,
	}
}

func (ur *userRepository) GetProfileInfo(uid int) (info models.JsonInfo, err error) {
	return
}

// func for testing
func (ur *userRepository) DeleteUser(mail string) error {
	return nil
}

func (ur *userRepository) UpdateUserPhotos(uid int, newImages *[]string) models.WorkMessage {
	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusNotImplemented),
		Status:  http.StatusNotImplemented,
	}
}

func (ur *userRepository) UpdUserGeneral(info models.JsonInfo, usr models.User) error {
	return nil
}

func (ur *userRepository) GetValidTags() ([]models.Tag, error) {
	return nil, nil
}

func (ur *userRepository) GetRandomFeedUsers(uid int, limit int, page int) ([]models.UserGeneral, error) {
	return nil, nil
}

func (ur *userRepository) GetUsersForChat(cid int64, users *models.UserGeneralList) models.WorkMessage {
	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusNotImplemented),
		Status:  http.StatusNotImplemented,
	}
}
