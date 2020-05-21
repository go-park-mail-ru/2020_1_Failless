package event

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	SearchEventsByUserPreferences(events *models.MidAndBigEventList, request *models.EventRequest) (int, error)
	CreateSmallEvent(smallEventForm *forms.SmallEventForm) (models.SmallEvent, error)
	UpdateSmallEvent(event *models.SmallEvent) (int, error)
	DeleteSmallEvent(uid int, eid int64) models.WorkMessage
	CreateMidEvent(midEventForm *forms.MidEventForm) (models.MidEvent, models.WorkMessage)
	JoinMidEvent(eventVote *models.EventFollow) models.WorkMessage
	LeaveMidEvent(eventVote *models.EventFollow) models.WorkMessage
	GetSmallEventsByUID(uid int64) (models.SmallEventList, error)
}
