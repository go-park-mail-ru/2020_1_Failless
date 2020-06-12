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
			Message: "Email is already in database",
			Status:  code,
		}
	}

	if err = eu.SendEmail(email); err != nil {
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

func (eu *emailUseCase) SendEmail(email *models.Email) error {
	to := []string{email.Email}
	var message []byte
	if email.Lang == "ru" {
		message = []byte(RU_Notify)
	} else if email.Lang == "es" {
		message = []byte(ES_Notify)
	} else {
		message = []byte(EN_Notify)
	}
    err := smtp.SendMail(settings.EmailServer.Host + ":" + settings.EmailServer.Port, settings.EmailServer.Auth, settings.EmailServer.Login, to, message)
	if err != nil {
		log.Println(err)
	}
    return err
}

const (
	RU_Notify = "Вы получили это сообщение, поскольку захотели, что бы мы уведомили Вас, когда сервис Eventum снова начнёт работать. Если это были не Вы - напишите нам"
	EN_Notify = "You've received this message since you asked us to notify you, when Eventum is back on its feet. If it weren't you - let us know"
	ES_Notify = "Has recibido este mensaje desde que nos has pedido que te lo notifiquemos cuando Eventum vuelva a ponerse de pie. Si no fuera usted - háganoslo saber"
)
