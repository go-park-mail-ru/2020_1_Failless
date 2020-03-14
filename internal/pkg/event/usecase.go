package event

import "failless/internal/pkg/models"

type UseCase interface {
	InitEventsByTime(events []models.Event) (int, error)
}
