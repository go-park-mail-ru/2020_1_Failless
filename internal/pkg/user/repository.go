package user

import "failless/internal/pkg/models"

// Storage interface, which works with user's data
// It can be in-memory realization or database realization
type Repository interface {
	GetUserByUID(uid int) (models.User, error)
	UpdateUserRating(uid int, rating float32) error
	SetUserLocation(uid int, point models.LocationPoint) error
	AddUserInfo(user models.User, info models.JsonInfo) error
	AddNewUser(user *models.User) error
	GetUserByPhoneOrEmail(phone string, email string) (models.User, error)
	GetProfileInfo(uid int) (models.JsonInfo, error)
	DeleteUser(mail string) error
}
