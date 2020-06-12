package usecase

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/email"
	"failless/internal/pkg/email/repository"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"log"
	"net/http"
	"net/smtp"
)

type emailUseCase struct {
	rep email.Repository
}

func GetUseCase() email.UseCase {
	return &emailUseCase{
		rep: repository.NewSqlEmailRepository(db.ConnectToDB()),
	}
}

func (eu *emailUseCase) SendReminder(email *models.Email) models.WorkMessage {
	code, err := eu.rep.SaveEmail(email)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  code,
		}
	}

	if err = eu.SendEmail([]string{email.Email}); err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: "",
			Status:  code,
		}
	}

	return models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusOK,
	}
}

func (eu *emailUseCase) SendEmail(to []string) error {
    message := []byte("This is a really unimaginative message, I know.")
    err := smtp.SendMail(settings.EmailServer.Host + ":" + settings.EmailServer.Port, settings.EmailServer.Auth, settings.EmailServer.Login, to, message)
	if err != nil {
		log.Println(err)
	}
    return err
}
