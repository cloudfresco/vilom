package common

import (
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

// MailerIntf interface to the Mailer
type MailerIntf interface {
	SendMail(msg Email) error
}

// MailerService Pointer to mailer
type MailerService struct {
	Mailer *gomail.Dialer
}

// Email - for sending email notifications
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
	Cc      string
}

// NewMailerService get connection to mailer and create a MailerService struct
func NewMailerService(mailerOpt *MailerOptions) (*MailerService, error) {

	mailer := gomail.NewDialer(mailerOpt.Server, mailerOpt.Port, mailerOpt.User, mailerOpt.Password)

	mailerService := &MailerService{}
	mailerService.Mailer = mailer

	return mailerService, nil
}

// SendMail - used for sending email
func (mailerService *MailerService) SendMail(msg Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailerService.Mailer.Username)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	err := mailerService.Mailer.DialAndSend(m)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 259}).Error(err)
		return err
	}
	return nil
}
