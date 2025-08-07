package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/services/notifier"
)

type NotificationServiceImpl struct {
	Dispatcher *notifier.Dispatcher
}

func NewNotificationService(dispatcher *notifier.Dispatcher) domain.NotificationService {
	return &NotificationServiceImpl{
		Dispatcher: dispatcher,
	}
}

func (s *NotificationServiceImpl) Send(notification notifier.Notification) error {
	return s.Dispatcher.NotifyAll(notification)
}
