package event

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	CreateEvent(event forms.EventForm) (models.Event, error)
	InitEventsByTime(events *[]models.Event) (int, error)
	InitEventsByKeyWords(events *[]models.Event, keys string, page int) (int, error)
	InitEventsByUserPreferences(events *[]models.Event, request *models.EventRequest) (int, error)
}
