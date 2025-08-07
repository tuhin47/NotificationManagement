package notifier

import (
	"NotificationManagement/config"
	"fmt"
	"net/smtp"
)

type EmailNotifier struct {
	smtp.Auth
	Address string
	From    string
}

func NewEmailNotifier() *EmailNotifier {
	e := config.Email()
	var auth smtp.Auth
	if e.Password != "" {
		auth = smtp.PlainAuth("", e.Username, e.Password, e.Host)

	}
	return &EmailNotifier{
		Auth:    auth,
		Address: fmt.Sprintf("%s:%d", e.Host, *e.Port),
		From:    e.From,
	}
}

func (e *EmailNotifier) Send(n Notification) error {
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", n.To, n.Subject, n.Message))
	return smtp.SendMail(e.Address, e.Auth, e.From, []string{n.To}, msg)
}

func (e *EmailNotifier) Type() string {
	return "email"
}
