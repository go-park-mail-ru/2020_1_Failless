package repository

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"net/http"
)

type eventsRepository struct {
}

func NewEventRepository() event.Repository {
	return &eventsRepository{}
}

func (er *eventsRepository) getEvents(withCondition string, sqlStatement string, args ...interface{}) ([]models.Event, error) {
	return nil, nil
}

func (er *eventsRepository) SaveNewEvent(event *models.Event) error {
	return nil
}

func (er *eventsRepository) GetNameByID(uid int) (string, error) {
	return "", nil
}

func (er *eventsRepository) GetAllEvents() ([]models.Event, error) {
	return nil, nil
}

func (er *eventsRepository) GetFeedEvents(uid int, limit int, page int) ([]models.Event, error) {
	return nil, nil
}

func (er *eventsRepository) GetEventsByKeyWord(keyWords string, page int) (models.EventList, error) {
	return nil, nil
}

func (er *eventsRepository) GetValidTags() ([]models.Tag, error) {
	return nil, nil
}

func (er *eventsRepository) GetNewEventsByTags(tags []int, uid int, limit int, page int) (models.EventList, error) {
	return nil, nil
}

func (er *eventsRepository) FollowMidEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) FollowBigEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) UnfollowMidEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) UnfollowBigEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) GetEventsWithFollowed(events *models.EventResponseList, request *models.EventRequest) error {
	return nil
}

func (er *eventsRepository) CreateSmallEvent(event *models.SmallEvent) error {
	return nil
}

func (er *eventsRepository) GetSmallEventsForUser(uid int) (models.SmallEventList, error) {
	return nil, nil
}

func (er *eventsRepository) UpdateSmallEvent(event *models.SmallEvent) (int, error) {
	return http.StatusNotImplemented, nil
}
