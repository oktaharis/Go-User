package email

import (
	"fmt"
	"net/smtp"
)

// EmailSender adalah sebuah struct yang merepresentasikan pengirim email.
type EmailSender struct {
	SMTPServer  string
	SMTPPort    string
	SenderEmail string
	Password    string
}

// NewEmailSender membuat sebuah instance EmailSender baru.
func NewEmailSender(smtpServer, smtpPort, senderEmail, password string) *EmailSender {
	return &EmailSender{
		SMTPServer:  smtpServer,
		SMTPPort:    smtpPort,
		SenderEmail: senderEmail,
		Password:    password,
	}
}

// SendEmail mengirimkan email dengan konten yang diberikan.
func (e *EmailSender) SendEmail(recipient, subject, body string) error {
	auth := smtp.PlainAuth("", e.SenderEmail, e.Password, e.SMTPServer)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", e.SenderEmail, recipient, subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", e.SMTPServer, e.SMTPPort), auth, e.SenderEmail, []string{recipient}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
