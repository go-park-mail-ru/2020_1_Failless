package event

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
)

type UseCase interface {
	InitEventsByTime(events []models.Event) (int, error)
	CreateEvent(event forms.EventForm) (models.Event, error)
	InitEventsByKeyWords(events *[]models.Event, keys string, page int) (int, error)
}
