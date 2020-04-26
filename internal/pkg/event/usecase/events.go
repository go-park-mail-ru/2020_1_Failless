package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/repository"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/settings"
	"log"
	"net/http"
)

type eventUseCase struct {
	rep event.Repository
}

func GetUseCase() event.UseCase {
	if settings.UseCaseConf.InHDD {
		log.Println("IN HDD")
		return &eventUseCase{
			rep: repository.NewSqlEventRepository(db.ConnectToDB()),
		}
	} else {
		log.Println("IN MEMORY")
		return &eventUseCase{
			rep: repository.NewEventRepository(),
		}
	}
}

func (ec *eventUseCase) InitEventsByTime(events *[]models.Event) (status int, err error) {
	*events, err = ec.rep.GetAllEvents()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (ec *eventUseCase) InitEventsByKeyWords(events *[]models.Event, keyWords string, page int) (status int, err error) {
	if keyWords == "" {
		*events, err = ec.rep.GetAllEvents()
	} else {
		*events, err = ec.rep.GetEventsByKeyWord(keyWords, page)
	}
	log.Println(events)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (ec *eventUseCase) CreateEvent(event forms.EventForm) (models.Event, error) {
	user, err := ec.rep.GetNameByID(event.UId)
	model := models.Event{}
	event.GetDBFormat(&model)
	model.Author = user
	err = ec.rep.SaveNewEvent(&model)
	return model, err
}

func (ec *eventUseCase) InitEventsByUserPreferences(events *[]models.Event, request *models.EventRequest) (int, error) {
	dbTags, err := ec.rep.GetValidTags()
	if err != nil {
		return http.StatusBadRequest, err
	}

	valid := ec.TakeValidTagsOnly(request.Tags, dbTags)
	log.Println(request)
	if valid != nil {
		*events, err = ec.rep.GetNewEventsByTags(valid, request.Uid, request.Limit, request.Page)
	} else {
		*events, err = ec.rep.GetFeedEvents(request.Uid, request.Limit, request.Page)
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	for i := 0; i < len(*events); i++ {
		(*events)[i].Tag = dbTags[(*events)[i].Type-1]
		log.Println(dbTags[(*events)[i].Type - 1])
	}

	log.Println(events)
	return http.StatusOK, nil
}

func (ec *eventUseCase) TakeValidTagsOnly(tagIds []int, tags []models.Tag) []int {
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

func (ec *eventUseCase) FollowEvent(subscription *models.EventFollow) network.Message {
	var err error
	if subscription.Type == "mid-event" {
		err = ec.rep.FollowMidEvent(subscription.Uid, subscription.Eid)
	} else if subscription.Type == "big-event" {
		err = ec.rep.FollowBigEvent(subscription.Uid, subscription.Eid)
	} else {
		return network.Message{
			Request: nil,
			Message: "Invalid subscription type",
			Status:  http.StatusBadRequest,
		}
	}

	if err != nil {
		return network.Message{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusConflict,
		}
	} else {
		return network.Message{
			Request: nil,
			Message: "OK",
			Status:  http.StatusCreated,
		}
	}
}
