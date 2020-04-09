package tag

import "failless/internal/pkg/models"

type Repository interface {
	GetAllTags() ([]models.Tag, error)
}
