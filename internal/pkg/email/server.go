package email

import (
	email2 "failless/configs/email"
	"log"
	"net/smtp"
	"os"
	"sync"
)

type serverCredentials struct {
	Login	string
	Pass	string
}

type SMTPServer struct {
	Creds	serverCredentials
	Host 	string
	Port 	string
	Auth	smtp.Auth
}

var server SMTPServer
var doOnce sync.Once

func ConnectToSMTPServer() {
	doOnce.Do(func() {
		server = SMTPServer{
			Host: "smtp.gmail.com",
			Port: "587",
			Creds: serverCredentials{
				Login: os.Getenv(email2.Secrets[0]),
				Pass:  os.Getenv(email2.Secrets[1]),
			},
		}
		server.Auth = smtp.PlainAuth("", server.Creds.Login, server.Creds.Pass, server.Host)
	})
}

func SendEmail(to []string) bool {
	ConnectToSMTPServer()
    message := []byte("This is a really unimaginative message, I know.")
    err := smtp.SendMail(server.Host + ":" + server.Port, server.Auth, server.Creds.Login, to, message)
    if err != nil {
        log.Println(err)
        return false
    } else {
		log.Println("Email Sent!")
		return true
	}
}
