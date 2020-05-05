package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/repository"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"fmt"
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

func (ec *eventUseCase) InitEventsByTime(events *models.EventList) (status int, err error) {
	*events, err = ec.rep.GetAllEvents()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (ec *eventUseCase) InitEventsByKeyWords(events *models.EventList, keys string, page int) (status int, err error) {
	if keys == "" {
		*events, err = ec.rep.GetAllEvents()
	} else {
		*events, err = ec.rep.GetEventsByKeyWord(keys, page)
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

func (ec *eventUseCase) CreateSmallEvent(event *models.SmallEvent) error {
	return ec.rep.CreateSmallEvent(event)
}

func (ec *eventUseCase) GetSmallEventsForUser(uid int) (models.SmallEventList, error) {
	return ec.rep.GetSmallEventsForUser(uid)
}


func (ec *eventUseCase) InitEventsByUserPreferences(events *models.EventList, request *models.EventRequest) (int, error) {
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
		log.Println(dbTags[(*events)[i].Type-1])
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

func (ec *eventUseCase) FollowEvent(subscription *models.EventFollow) models.WorkMessage {
	var err error
	if subscription.Type == "mid-event" {
		err = ec.rep.FollowMidEvent(subscription.Uid, subscription.Eid)
	} else if subscription.Type == "big-event" {
		err = ec.rep.FollowBigEvent(subscription.Uid, subscription.Eid)
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "Invalid subscription type",
			Status:  http.StatusBadRequest,
		}
	}

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusConflict,
		}
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "OK",
			Status:  http.StatusCreated,
		}
	}
}

func (ec *eventUseCase) UnfollowEvent(subscription *models.EventFollow) models.WorkMessage {
	var err error
	if subscription.Type == "mid-event" {
		err = ec.rep.UnfollowMidEvent(subscription.Uid, subscription.Eid)
	} else if subscription.Type == "big-event" {
		err = ec.rep.UnfollowBigEvent(subscription.Uid, subscription.Eid)
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "Invalid subscription type",
			Status:  http.StatusBadRequest,
		}
	}

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusConflict,
		}
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "OK",
			Status:  http.StatusCreated,
		}
	}
}

func (ec *eventUseCase) SearchEventsByUserPreferences(events *models.EventResponseList, request *models.EventRequest) (int, error) {
	var err error
	if request.Uid == 0 {
		tempEvents, _ := ec.rep.GetAllEvents()
		for _, tempEvent := range tempEvents {
			tempEventResponse := models.EventResponse{
				Event:    tempEvent,
				Followed: false,
			}
			*events = append(*events, tempEventResponse)
		}
	} else {
		err = ec.rep.GetEventsWithFollowed(events, request)
		if err != nil {
			fmt.Println(err)
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func (ec *eventUseCase) UpdateSmallEvent(event *models.SmallEvent) (int, error) {
	return ec.rep.UpdateSmallEvent(event)
}

func (ec *eventUseCase) DeleteSmallEvent(uid int, eid int64) models.WorkMessage {
	err := ec.rep.DeleteSmallEvent(uid, eid)

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusOK),
		Status:  http.StatusOK,
	}
}
