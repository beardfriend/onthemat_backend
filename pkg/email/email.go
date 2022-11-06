package email

import (
	"net/smtp"

	"onthemat/internal/app/config"
)

type Email struct {
	auth     smtp.Auth
	username string
	host     string
}

func NewEmail(config *config.Config) *Email {
	auth := smtp.PlainAuth("", config.Email.UserName, config.Email.Password, config.Email.Host)
	return &Email{
		auth:     auth,
		username: config.Email.UserName,
		host:     config.Email.Host,
	}
}

func (e *Email) Send(to []string, msg []byte) error {
	return smtp.SendMail(e.host+":587", e.auth, e.username, to, msg)
}
