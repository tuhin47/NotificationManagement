package notifier

import "fmt"

type TelegramNotifier struct{}

func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{}
}

func (t *TelegramNotifier) Send(n Notification) error {
	fmt.Printf("[Telegram] To: %s, Message: %s\n", n.To, n.Message)
	return nil
}

func (t *TelegramNotifier) Type() string {
	return "telegram"
}
