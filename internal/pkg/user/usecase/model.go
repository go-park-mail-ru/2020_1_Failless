package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/repository"
	"log"
	"net/http"
)

type UserUseCase struct {
	Rep user.Repository
}

func GetUseCase() user.UseCase {
	if settings.UseCaseConf.InHDD {
		log.Println("IN HDD")
		return &UserUseCase{
			Rep: repository.NewSqlUserRepository(db.ConnectToDB()),
		}
	} else {
		log.Println("IN MEMORY")
		return &UserUseCase{
			Rep: repository.NewUserRepository(),
		}
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

func (uc *UserUseCase) TakeValidTagsOnly(tagIds []int, tags []models.Tag) []int {
	var valid []int = nil
	for _, tagId := range tagIds {
		for _, tag := range tags {
			if tagId == tag.TagId {
				valid = append(valid, tagId)
			}
		}
	}

	return valid
}

func (uc *UserUseCase) GetUserSubscriptions(events *models.EventList, uid int) (int, error) {
	var err error
	*events, err = uc.Rep.GetUserSubscriptions(uid)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (uc *UserUseCase) GetFeedResults(
	users *[]models.UserGeneral,
	form *[]forms.GeneralForm) (models.FeedResults, error) {
	// Get subscriptions for FeedUsers
	var events [][]models.Event
	for i := 0; i < len(*users); i++ {
		var userEvents models.EventList
		if _, err := uc.GetUserSubscriptions(&userEvents, (*users)[i].Uid); err != nil {
			return nil, err
		}
		events = append(events, userEvents)
	}

	// Collecting all together
	var result models.FeedResults
	for i := 0; i < len(*users); i++ {
		post := models.FeedPost{}
		post.Uid = (*users)[i].Uid
		post.Name = (*users)[i].Name
		post.Photos = (*users)[i].Photos
		post.About = (*users)[i].About
		post.Birthday = (*users)[i].Birthday
		post.Gender = (*users)[i].Gender
		post.Events = (*form)[i].Events
		post.Tags = (*form)[i].Tags
		if events != nil {
			post.Subs = events[i]
		}
		result = append(result, post)
	}

	return result, nil
}
