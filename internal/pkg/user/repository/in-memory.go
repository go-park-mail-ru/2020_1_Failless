package repository

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
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
func (ur *userRepository) getUser(sqlStatement string, args ...interface{}) (user models.User, err error) {
	return models.User{}, nil
}

// Private method
func (ur *userRepository) getEvents(sqlStatement string, args ...interface{}) ([]models.Event, error) {
	return nil, nil
}

func (ur *userRepository) AddNewUser(user *models.User) error {
	return nil
}

func (ur *userRepository) AddUserInfo(credentials models.User, info models.JsonInfo) error {
	return nil
}

func (ur *userRepository) SetUserLocation(uid int, point models.LocationPoint) error {
	return nil
}

func (ur *userRepository) UpdateUserRating(uid int, rating float32) error {
	return nil
}

func (ur *userRepository) UpdateUserTags(uid int, tagId int) error {
	return nil
}

func (ur *userRepository) UpdateUserSimple(uid int, social []string, about *string) error {
	return nil
}

func (ur *userRepository) GetProfileInfo(uid int) (info models.JsonInfo, err error) {
	return
}

// func for testing
func (ur *userRepository) DeleteUser(mail string) error {
	return nil
}

func (ur *userRepository) GetUserEvents(uid int) ([]models.Event, error) {
	return nil, nil
}

func (ur *userRepository) GetEventsByTag(tag string) ([]models.Event, error) {
	return nil, nil
}

func (ur *userRepository) GetUserTags(uid int) ([]models.Tag, error) {
	return nil, nil
}

func (ur *userRepository) UpdateUserPhotos(uid int, name string) error {
	return nil
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
