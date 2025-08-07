package domain

import (
	"NotificationManagement/services/notifier"
)

type NotificationService interface {
	Send(notification notifier.Notification) error
}
