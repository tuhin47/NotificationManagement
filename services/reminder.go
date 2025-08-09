package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"context"
)

type ReminderServiceImpl struct {
	domain.CommonService[models.Reminder]
	domain.NotificationDispatcher
}

func NewReminderService(repo domain.ReminderRepository, dispatcher domain.NotificationDispatcher) domain.ReminderService {
	service := &ReminderServiceImpl{
		NotificationDispatcher: dispatcher,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}
func (r ReminderServiceImpl) SendReminders(ctx context.Context, reminderId uint) error {

	reminder, err := r.CommonService.GetModelById(ctx, reminderId, &[]string{"Request", "Request.User", "Request.User.Telegram"})
	if err != nil {
		return err
	}
	return r.NotificationDispatcher.Notify(ctx, &types.Notification{
		// TODO : Need to collect from request/reminder
		Subject:  reminder.Message,
		Message:  reminder.Message,
		Channels: []string{"sms", "email", "telegram"},
		UserId:   reminder.Request.UserID,
		User:     reminder.Request.User,
	})
}
