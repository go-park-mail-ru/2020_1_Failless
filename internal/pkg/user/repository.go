package user

import (
	"failless/internal/pkg/models"
)

// Storage interface, which works with user's data
// It can be in-memory realization or database realization
type Repository interface {
	GetUserByUID(uid int) (models.User, error)
	UpdateUserTags(uid int, tagId int) error
	UpdateUserSimple(uid int, social []string, about *string) error
	UpdateUserRating(uid int, rating float32) error
	UpdateUserPhotos(uid int, name string) error
	SetUserLocation(uid int, point models.LocationPoint) error
	UpdUserGeneral(info models.JsonInfo, user models.User) error
	GetValidTags() ([]models.Tag, error)
	GetRandomFeedUsers(uid int, limit int, page int) ([]models.UserGeneral, error)

	// Deprecated: use UpdateUserTags, UpdateUserSimple,
	// UpdateUserPhotos, UpdateUserRating instead
	AddUserInfo(user models.User, info models.JsonInfo) error
	AddNewUser(user *models.User) error
	GetUserByPhoneOrEmail(phone string, email string) (models.User, error)
	GetProfileInfo(uid int) (models.JsonInfo, error)
	GetUserEvents(uid int) ([]models.Event, error)
	GetUserTags(uid int) ([]models.Tag, error)
	DeleteUser(mail string) error
}
