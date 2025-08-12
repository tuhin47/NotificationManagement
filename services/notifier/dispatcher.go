package notifier

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"context"
)

type Dispatcher struct {
	Notifiers *[]domain.Notifier
	domain.UserService
}

func NewNotificationDispatcher(email *EmailNotifier, sms *SMSNotifier, telegram domain.TelegramNotifier, service domain.UserService) domain.NotificationDispatcher {
	return &Dispatcher{
		Notifiers:   &[]domain.Notifier{email, sms, telegram},
		UserService: service,
	}
}

func (d *Dispatcher) Notify(ctx context.Context, notification *types.Notification) error {
	if notification.User == nil {
		user, err := d.UserService.GetModelById(ctx, notification.UserId, &[]string{"Telegram"})
		if err != nil {
			return err
		}
		notification.User = user
	}
	for _, notifier := range *d.GetDispatchers() {
		for _, channel := range notification.Channels {
			if notifier.Type() == channel && notifier.IsActive() {
				if err := notifier.Send(ctx, notification); err != nil {
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
