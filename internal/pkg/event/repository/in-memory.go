package repository

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
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

func (er *eventsRepository) GetEventsByKeyWord(keyWordsString string, page int) ([]models.Event, error) {
	return nil, nil
}

func (er *eventsRepository) GetValidTags() ([]models.Tag, error) {
	return nil, nil
}

func (er *eventsRepository) GetNewEventsByTags(tags []int, uid int, limit int, page int) ([]models.Event, error) {
	return nil, nil
}

func (er *eventsRepository) FollowMidEvent(uid, eid int) error {
	return nil
}
func (er *eventsRepository) FollowBigEvent(uid, eid int) error {
	return nil
}
