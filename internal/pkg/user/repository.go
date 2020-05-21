package user

import (
	"failless/internal/pkg/models"
)

// Storage interface, which works with user's data
// It can be in-memory realization or database realization
type Repository interface {
	GetUserByUID(uid int) (models.User, error)
	UpdateUserTags(uid int, tagIDs []int) models.WorkMessage
	UpdateUserRating(uid int, rating float32) error
	UpdateUserPhotos(uid int, newImages *[]string) models.WorkMessage
	UpdateUserAbout(uid int, about string) models.WorkMessage
	SetUserLocation(uid int, point models.LocationPoint) error
	UpdUserGeneral(info models.JsonInfo, user models.User) error

	// Deprecated: use UpdateUserTags, UpdateUserSimple,
	// UpdateUserPhotos, UpdateUserRating instead
	AddNewUser(user *models.User) error
	GetUserByPhoneOrEmail(phone string, email string) (models.User, error)
	GetProfileInfo(uid int) (models.JsonInfo, error)
	DeleteUser(mail string) error
	GetValidTags() ([]models.Tag, error)
	GetRandomFeedUsers(uid int, limit int, page int) ([]models.UserGeneral, error)
	GetUsersForChat(cid int64, users *models.UserGeneralList) models.WorkMessage
}
