package usecase

//go:generate mockgen -destination=../mocks/mock_usecase.go -package=mocks failless/internal/pkg/user UseCase

import (
	"errors"
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	eventRep "failless/internal/pkg/event/repository"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
	"log"
	"net/http"
)

const (
	MessageUserDoesntExist = "User doesn't exist\n"
)

var (
	CorrectMessage = models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusOK,
	}
)

type UserUseCase struct {
	Rep user.Repository
	eventRep event.Repository
}

func GetUseCase() user.UseCase {
	return &UserUseCase{
		Rep: repository.NewSqlUserRepository(db.ConnectToDB()),
		eventRep: eventRep.NewSqlEventRepository(db.ConnectToDB()),
	}
}

func (uc *UserUseCase) InitUsersByUserPreferences(users *[]models.UserGeneral, request *models.UserRequest) (int, error) {
	_, err := uc.Rep.GetValidTags()
	if err != nil {
		return http.StatusBadRequest, err
	}

	*users, err = uc.Rep.GetRandomFeedUsers(request.Uid, request.Limit, request.Page)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	log.Println(users)
	return http.StatusOK, nil
}

func (uc *UserUseCase) GetUserSubscriptions(subscriptions *models.MidAndBigEventList, uid int) models.WorkMessage {
	code, err := uc.eventRep.GetSubscriptionMidEvents(&subscriptions.MidEvents, uid)
	//code, err != er.GetSubscriptionBigEvents(&subscriptions.BigEvents, uid)

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}
	return CorrectMessage
}

func (uc *UserUseCase) GetFeedResultsFor(uid int, users *[]models.UserGeneral) (models.FeedResults, models.WorkMessage) {
	var feedResults models.FeedResults

	for _, feedUser := range *users {
		var feedResult models.FeedPost
		feedResult.Uid = feedUser.Uid
		feedResult.Photos = feedUser.Photos
		feedResult.TagsId = feedUser.TagsId
		feedResult.About = feedUser.About
		feedResult.Name = feedUser.Name
		feedResult.Birthday = feedUser.Birthday
		feedResult.Gender = feedUser.Gender

		code, err := uc.eventRep.GetSmallEventsForUser(&feedResult.OnwEvents.SmallEvents, feedResult.Uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}
		code, err = uc.eventRep.GetOwnMidEventsWithAnotherUserFollowed(&feedResult.OnwEvents.MidEvents, feedResult.Uid, uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}
		code, err = uc.eventRep.GetSubscriptionMidEventsWithAnotherUserFollowed(&feedResult.Subscriptions.MidEvents, feedResult.Uid, uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}

		feedResults = append(feedResults, feedResult)
	}

	return feedResults, CorrectMessage
}

func (uc *UserUseCase) UpdateUserAbout(uid int, about string) models.WorkMessage {
	return uc.Rep.UpdateUserAbout(uid, about)
}

func (uc *UserUseCase) UpdateUserTags(uid int, tagIDs []int) models.WorkMessage {
	return uc.Rep.UpdateUserTags(uid, tagIDs)
}

func (uc *UserUseCase) GetUserInfo(profile *forms.GeneralForm) (int, error) {
	row, err := uc.Rep.GetProfileInfo(profile.Uid)
	if err != nil {
		log.Println(err.Error())
		return http.StatusNotFound, err
	}

	profile.FillProfile(row)

	base, err := uc.Rep.GetUserByUID(profile.Uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	(*profile).SignForm.Name = base.Name
	(*profile).Phone = base.Phone
	(*profile).Email = base.Email
	(*profile).Uid = base.Uid

	return http.StatusOK, nil
}

func (uc *UserUseCase) UpdateUserPhotos(uid int, newImages *forms.EImageList) models.WorkMessage {
	var imageNames []string
	for index := range *newImages {
		imageNames = append(imageNames, (*newImages)[index].ImgName)
		(*newImages)[index].ImgBase64 = ""
	}
	return uc.Rep.UpdateUserPhotos(uid, &imageNames)
}

func (uc *UserUseCase) UpdateUserBase(form *forms.SignForm) (int, error) {
	usr, err := uc.Rep.GetUserByUID(form.Uid)
	if err != nil {
		log.Println(err.Error())
		return http.StatusNotFound, err
	}

	usr.Uid = form.Uid
	//usr.Name = form.Name
	usr.Email = form.Email
	usr.Phone = form.Phone
	usr.Password, err = security.EncryptPassword(form.Password)
	if err != nil {
		log.Println(err.Error())
	}

	var inf = models.JsonInfo{}
	//inf.Birthday = form.Birthday
	//inf.Gender = form.Gender

	if err := uc.Rep.UpdUserGeneral(inf, usr); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (uc *UserUseCase) GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) models.WorkMessage {
	code, err := uc.eventRep.GetSmallEventsForUser(smallEvents, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}
	return CorrectMessage
}

func (uc *UserUseCase) GetUserOwnEvents(ownEvents *models.OwnEventsList, uid int) models.WorkMessage {
	// Get Mid events
	midEventList := models.MidEventList{}
	code, err := uc.eventRep.GetOwnMidEvents(&midEventList, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	// Get Small events
	smallEventList := models.SmallEventList{}
	code, err = uc.eventRep.GetSmallEventsForUser(&smallEventList, uid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	// Assign and return
	ownEvents.MidEvents = midEventList
	ownEvents.SmallEvents = smallEventList
	return CorrectMessage
	// TODO: rewrite in goroutines
}

func (uc *UserUseCase) FillFormIfExist(cred *models.User) (int, error) {
	log.Println(*cred)
	user, err := uc.Rep.GetUserByPhoneOrEmail(cred.Phone, cred.Email)
	if err == nil && user.Uid < 0 {
		log.Println("user not found")
		return http.StatusNotFound, errors.New(MessageUserDoesntExist)
	} else if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}

	*cred = user
	return http.StatusOK, nil
}

func (uc *UserUseCase) RegisterNewUser(user *forms.SignForm) error {
	// TODO: move it to repository
	bPass, err := security.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	dbUser := models.User{
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: bPass,
	}

	return uc.Rep.AddNewUser(&dbUser)
}
