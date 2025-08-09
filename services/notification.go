package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
)

type NotificationServiceImpl struct {
	Dispatcher domain.NotificationDispatcher
}

func NewNotificationService(dispatcher domain.NotificationDispatcher) domain.NotificationService {
	return &NotificationServiceImpl{
		Dispatcher: dispatcher,
	}
}

func (s *NotificationServiceImpl) Send(notification *types.Notification) error {
	return s.Dispatcher.Notify(notification)
}
