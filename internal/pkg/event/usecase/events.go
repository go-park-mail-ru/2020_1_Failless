package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/repository"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"log"
	"net/http"
)

type eventUseCase struct {
	rep event.Repository
}

func GetUseCase() event.UseCase {
	return &eventUseCase{
		rep: repository.NewSqlEventRepository(db.ConnectToDB()),
	}
}

func (ec *eventUseCase) GetSmallEventsByUID(uid int64) (models.SmallEventList, error) {
	// TODO: implement it
	return nil, nil
}

func (ec *eventUseCase) CreateSmallEvent(smallEventForm *forms.SmallEventForm) (models.SmallEvent, error) {
	smallEvent := models.SmallEvent{}
	smallEventForm.GetDBFormat(&smallEvent)
	err := ec.rep.CreateSmallEvent(&smallEvent)
	return smallEvent, err
}

func (ec *eventUseCase) CreateMidEvent(midEventForm *forms.MidEventForm) (models.MidEvent, models.WorkMessage) {
	midEvent := models.MidEvent{}
	midEventForm.GetDBFormat(&midEvent)

	err := ec.rep.CreateMidEvent(&midEvent)

	if err != nil {
		return midEvent, models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return midEvent, models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusCreated,
	}
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

func (ec *eventUseCase) SearchEventsByUserPreferences(events *models.MidAndBigEventList, request *models.EventRequest) (int, error) {
	if request.Uid == 0 {
		code, err := ec.rep.GetAllMidEvents(&events.MidEvents, request)
		//bigEvents, _ := ec.rep.GetAllBigEvents()
		if err != nil {
			log.Println(err)
			return code, err
		}
	} else {
		code, err := ec.rep.GetMidEventsWithFollowed(&events.MidEvents, request)
		if err != nil {
			log.Println(err)
			return code, err
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

func (ec *eventUseCase) JoinMidEvent(eventVote *models.EventFollow) models.WorkMessage {
	code, err := ec.rep.JoinMidEvent(eventVote.Uid, eventVote.Eid)

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusOK),
		Status:  http.StatusOK,
	}
}

func (ec *eventUseCase) LeaveMidEvent(eventVote *models.EventFollow) models.WorkMessage {
	code, err := ec.rep.LeaveMidEvent(eventVote.Uid, eventVote.Eid)

	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: http.StatusText(http.StatusOK),
		Status:  http.StatusOK,
	}
}
