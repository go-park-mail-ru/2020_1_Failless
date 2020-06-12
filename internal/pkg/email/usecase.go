package email

import "failless/internal/pkg/models"

type UseCase interface {
	SendReminder(email *models.Email) models.WorkMessage

	SendEmail(email *models.Email) error
}
