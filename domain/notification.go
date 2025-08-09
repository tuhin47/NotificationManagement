package domain

import (
	"NotificationManagement/types"
)

type Notifier interface {
	Send(*types.Notification) error
	Type() string
	IsActive() bool
}

type NotificationDispatcher interface {
	GetDispatchers() *[]Notifier
}

type NotificationService interface {
	Send(notification *types.Notification) error
}
