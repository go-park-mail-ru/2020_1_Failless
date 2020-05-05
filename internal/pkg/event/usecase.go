package event

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	CreateEvent(event forms.EventForm) (models.Event, error)
	InitEventsByTime(events *models.EventList) (int, error)
	InitEventsByKeyWords(events *models.EventList, keys string, page int) (int, error)
	InitEventsByUserPreferences(events *models.EventList, request *models.EventRequest) (int, error)
	FollowEvent(subscription *models.EventFollow) models.WorkMessage
	UnfollowEvent(subscription *models.EventFollow) models.WorkMessage
	SearchEventsByUserPreferences(events *models.EventResponseList, request *models.EventRequest) (int, error)
	CreateSmallEvent(event *models.SmallEvent) error
	UpdateSmallEvent(event *models.SmallEvent) (int, error)
	DeleteSmallEvent(uid int, eid int64) models.WorkMessage
	GetSmallEventsForUser(uid int) (models.SmallEventList, error)
}
