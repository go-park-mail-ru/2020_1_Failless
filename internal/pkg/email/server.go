package email

import (
	"failless/internal/pkg/settings"
	"log"
	"net/smtp"
)

func SendEmail(to []string) bool {
    message := []byte("This is a really unimaginative message, I know.")
    err := smtp.SendMail(settings.EmailServer.Host + ":" + settings.EmailServer.Port, settings.EmailServer.Auth, settings.EmailServer.Login, to, message)
    if err != nil {
        log.Println(err)
        return false
    } else {
		log.Println("Email Sent!")
		return true
	}
}
