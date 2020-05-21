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

//func (er *eventsRepository) getEvents(withCondition string, sqlStatement string, args ...interface{}) ([]models.Event, error) {
//	return nil, nil
//}

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

func (er *eventsRepository) FollowBigEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) UnfollowBigEvent(uid, eid int) error {
	return nil
}

func (er *eventsRepository) CreateSmallEvent(event *models.SmallEvent) error {
	return nil
}

func (er *eventsRepository) GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) UpdateSmallEvent(event *models.SmallEvent) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) DeleteSmallEvent(uid int, eid int64) error {
	return nil
}

func (er *eventsRepository) CreateMidEvent(event *models.MidEvent) error {
	return nil
}

func (er *eventsRepository) GetOwnEventsForUser(ownEvents *models.OwnEventsList, uid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetOwnMidEvents(midEvents *models.MidEventList, uid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetMidEventsWithFollowed(midEvents *models.MidEventList, request *models.EventRequest) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetAllMidEvents(midEvents *models.MidEventList, request *models.EventRequest) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) JoinMidEvent(uid, eid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) LeaveMidEvent(uid, eid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetSubscriptionMidEvents(midEvent *models.MidEventList, uid int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetOwnMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, admin, member int) (int, error) {
	return http.StatusNotImplemented, nil
}

func (er *eventsRepository) GetSubscriptionMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, uid, visitor int) (int, error) {
	return http.StatusNotImplemented, nil
}
