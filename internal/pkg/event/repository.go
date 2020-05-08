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
	JoinMidEvent(uid, eid int) (int, error)
	FollowBigEvent(uid, eid int) error
	UnfollowMidEvent(uid, eid int) error
	UnfollowBigEvent(uid, eid int) error
	CreateSmallEvent(event *models.SmallEvent) error
	UpdateSmallEvent(event *models.SmallEvent) (int, error)
	DeleteSmallEvent(uid int, eid int64) error
	GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) (int, error)
	CreateMidEvent(event *models.MidEvent) error
	GetMidEventsForUser(midEvents *models.MidEventList, uid int) (int, error)
	GetAllMidEvents(midEvents *models.MidEventList, request *models.EventRequest) (int, error)
	GetMidEventsWithFollowed(midEvents *models.MidEventList, request *models.EventRequest) (int, error)
}
