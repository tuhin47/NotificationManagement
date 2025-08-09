package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/services/notifier"
	"NotificationManagement/types"
)

type NotificationServiceImpl struct {
	Dispatcher *notifier.Dispatcher
}

func NewNotificationService(dispatcher *notifier.Dispatcher) domain.NotificationService {
	return &NotificationServiceImpl{
		Dispatcher: dispatcher,
	}
}

func (s *NotificationServiceImpl) Send(notification *types.Notification) error {
	return s.Dispatcher.Notify(notification)
}
