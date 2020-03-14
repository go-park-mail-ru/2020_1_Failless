package tag

import "failless/internal/pkg/models"

type UseCase interface {
	InitEventsByTime(tags []models.Tag) (int, error)
}
