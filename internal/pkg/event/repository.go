package event

import "failless/internal/pkg/models"

type Repository interface {
	GetAllEvents() ([]models.Event, error)
	SaveNewEvent(event *models.Event) error
	GetNameByID(uid int) (string, error)
	GetFeedEvents(limit int, page int) ([]models.Event, error)
	GetEventsByKeyWord(keyWords string, page int) ([]models.Event, error)
	GetValidTags([]int) ([]int, error)
	GetNewEventsByTags(tags []int, uid int, limit int, page int) ([]models.Event, error)
}
