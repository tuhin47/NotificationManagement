package notifier

import (
	"NotificationManagement/logger"
	"NotificationManagement/types"
	"context"
	"fmt"
)

type SMSNotifier struct{}

func NewSMSNotifier() *SMSNotifier {
	return &SMSNotifier{}
}

func (s *SMSNotifier) Send(ctx context.Context, notification *types.Notification) error {
	logger.Info(fmt.Sprintf("SMS =>To: %s, Message: %s", notification.UserId, notification.Message))
	return nil
}

func (s *SMSNotifier) Type() string {
	return "sms"
}

func (s *SMSNotifier) IsActive() bool {
	return true
}
