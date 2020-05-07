package user

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	UpdateUserMeta(form *forms.MetaForm) (int, error)
	UpdateUserInfo(form *forms.GeneralForm) (int, error)
	GetUserInfo(profile *forms.GeneralForm) (int, error)
	FillFormIfExist(cred *models.User) (int, error)
	RegisterNewUser(user *forms.SignForm) error
	AddImageToProfile(uid int, name string) error
	UpdateUserBase(form *forms.SignForm) (int, error)
	InitUsersByUserPreferences(users *[]models.UserGeneral, request *models.UserRequest) (int, error)
	GetUserSubscriptions(events *models.EventList, uid int) (int, error)
	GetFeedResults(users *[]models.UserGeneral, form *[]forms.GeneralForm) (models.FeedResults, error)
	GetUserOwnEvents(ownEvents *models.OwnEventsList, uid int) models.WorkMessage
	GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) models.WorkMessage
}

// Get gender string by int id
func GenderById(genderId int) string {
	switch genderId {
	case models.Male:
		return "male"
	case models.Female:
		return "female"
	}
	return "other"
}

// Get gender id by string name
func GenderByStr(gender string) int {
	switch gender {
	case "male":
		return models.Male
	case "female":
		return models.Female
	}

	// Other gender
	return models.Other
}
