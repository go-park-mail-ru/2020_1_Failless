package event

import "failless/internal/pkg/models"

type Repository interface {
	GetAllEvents() ([]models.Event, error)
}
