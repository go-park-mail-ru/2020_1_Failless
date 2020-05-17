package usecase

import (
	"failless/internal/pkg/db"
	eventRep "failless/internal/pkg/event/repository"
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

func (uc *UserUseCase) GetUserSubscriptions(subscriptions *models.MidAndBigEventList, uid int) models.WorkMessage {
	er := eventRep.NewSqlEventRepository(db.ConnectToDB())

	code, err := er.GetSubscriptionMidEvents(&subscriptions.MidEvents, uid)
	//code, err != er.GetSubscriptionBigEvents(&subscriptions.BigEvents, uid)

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "",
			Status:  code,
		}
	}
}

func (uc *UserUseCase) GetFeedResultsFor(uid int, users *[]models.UserGeneral) (models.FeedResults, models.WorkMessage) {
	er := eventRep.NewSqlEventRepository(db.ConnectToDB())

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

		code, err := er.GetSmallEventsForUser(&feedResult.OnwEvents.SmallEvents, feedResult.Uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}
		code, err = er.GetOwnMidEventsWithAnotherUserFollowed(&feedResult.OnwEvents.MidEvents, feedResult.Uid, uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}
		code, err = er.GetSubscriptionMidEventsWithAnotherUserFollowed(&feedResult.Subscriptions.MidEvents, feedResult.Uid, uid)
		if err != nil {
			return feedResults, models.WorkMessage{
				Request: nil,
				Message: err.Error(),
				Status:  code,
			}
		}

		feedResults = append(feedResults, feedResult)
	}

	return feedResults, models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusOK,
	}
}
