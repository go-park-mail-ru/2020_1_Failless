package event

import "failless/internal/pkg/models"

type Repository interface {
	GetAllEvents() ([]models.Event, error)
	SaveNewEvent(event *models.Event) error
	GetNameByID(uid int) (string, error)
	GetFeedEvents(uid int, limit int, page int) ([]models.Event, error)
	GetEventsByKeyWord(keyWords string, page int) (models.EventList, error)
	GetValidTags() ([]models.Tag, error)
	GetNewEventsByTags(tags []int, uid int, limit int, page int) (models.EventList, error)
	FollowMidEvent(uid, eid int) error
	FollowBigEvent(uid, eid int) error
	GetEventsWithFollowed(events *models.EventResponseList, request *models.EventRequest) error
	CreateSmallEvent(event *models.SmallEvent) error
	GetSmallEventsForUser(uid int) (models.SmallEventList, error)
}
