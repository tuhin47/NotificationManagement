package notifier

import (
	"NotificationManagement/logger"
	"fmt"
)

type TelegramNotifier struct{}

func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{}
}

func (t *TelegramNotifier) Send(n Notification) error {
	logger.Info(fmt.Sprintf("[Telegram] To: %s, Message: %s", n.To, n.Message), n)
	return nil
}

func (t *TelegramNotifier) Type() string {
	return "telegram"
}
