package email

import "failless/internal/pkg/models"

type Repository interface {
	SaveEmail(email *models.Email) (int, error)
}
