package notifier

import (
	"NotificationManagement/logger"
	"NotificationManagement/types"
	"fmt"
)

type SMSNotifier struct{}

func NewSMSNotifier() *SMSNotifier {
	return &SMSNotifier{}
}

func (s *SMSNotifier) Send(n *types.Notification) error {
	logger.Info(fmt.Sprintf("SMS =>To: %s, Message: %s", n.To, n.Message))
	return nil
}

func (s *SMSNotifier) Type() string {
	return "sms"
}

func (s *SMSNotifier) IsActive() bool {
	return true
}
