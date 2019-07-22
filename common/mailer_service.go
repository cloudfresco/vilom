package common

import (
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

// MailerIntf interface to the Mailer
type MailerIntf interface {
	SendMail(msg Email, gomailer *gomail.Dialer) error
}

// Mailer Pointer to mailer
type Mailer struct {
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

// SendMail - used for sending email
func SendMail(msg Email, gomailer *gomail.Dialer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", gomailer.Username)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	err := gomailer.DialAndSend(m)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 259}).Error(err)
		return err
	}
	return nil
}
