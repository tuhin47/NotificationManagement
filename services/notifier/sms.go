package notifier

import "fmt"

type SMSNotifier struct{}

func NewSMSNotifier() *SMSNotifier {
	return &SMSNotifier{}
}

func (s *SMSNotifier) Send(n Notification) error {
	fmt.Printf("[SMS] To: %s, Message: %s\n", n.To, n.Message)
	return nil
}

func (s *SMSNotifier) Type() string {
	return "sms"
}
