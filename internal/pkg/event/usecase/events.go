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

type userUseCase struct {
	rep event.Repository
}

func GetUseCase() event.UseCase {
	return &userUseCase{
		rep: repository.NewSqlEventRepository(db.ConnectToDB()),
	}
}

func (uc *userUseCase) InitEventsByTime(events *[]models.Event) (status int, err error) {
	*events, err = uc.rep.GetAllEvents()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (uc *userUseCase) InitEventsByKeyWords(events *[]models.Event, keyWords string, page int) (status int, err error) {
	*events, err = uc.rep.GetEventsByKeyWord(keyWords, page)
	log.Println(events)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (uc *userUseCase) CreateEvent(event forms.EventForm) (models.Event, error) {
	user, err := uc.rep.GetNameByID(event.UId)
	model := models.Event{}
	event.GetDBFormat(&model)
	model.Author = user
	err = uc.rep.SaveNewEvent(&model)
	return model, err
}

func (uc *userUseCase) InitEventsByUserPreferences(events *[]models.Event, request *models.EventRequest) (int, error) {
	dbTags, err := uc.rep.GetValidTags()
	if err != nil {
		return http.StatusBadRequest, err
	}

	valid := uc.TakeValidTagsOnly(request.Tags, dbTags)
	log.Println(request)
	if valid != nil {
		*events, err = uc.rep.GetNewEventsByTags(valid, request.Uid, request.Limit, request.Page)
	} else {
		*events, err = uc.rep.GetFeedEvents(request.Uid, request.Limit, request.Page)
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

func (uc *userUseCase) TakeValidTagsOnly(tagIds []int, tags []models.Tag) []int {
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
