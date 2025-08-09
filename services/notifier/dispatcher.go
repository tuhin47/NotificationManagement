package notifier

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
)

type Dispatcher struct {
	Notifiers *[]domain.Notifier
}

func NewNotificationDispatcher(email *EmailNotifier, sms *SMSNotifier, telegram domain.TelegramNotifier) domain.NotificationDispatcher {
	return &Dispatcher{
		Notifiers: &[]domain.Notifier{email, sms, telegram},
	}
}

func (d *Dispatcher) Notify(n *types.Notification) error {
	for _, notifier := range *d.GetDispatchers() {
		for _, channel := range n.Channels {
			if notifier.Type() == channel && notifier.IsActive() {
				if err := notifier.Send(n); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *Dispatcher) GetDispatchers() *[]domain.Notifier {
	return d.Notifiers
}
