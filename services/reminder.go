package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"context"
	"fmt"
)

type ReminderServiceImpl struct {
	domain.CommonService[models.Reminder]
	domain.NotificationDispatcher
	domain.AiDispatcher
}

func NewReminderService(repo domain.ReminderRepository, dispatcher domain.NotificationDispatcher, aiDispatcher domain.AiDispatcher) domain.ReminderService {
	service := &ReminderServiceImpl{
		NotificationDispatcher: dispatcher,
		AiDispatcher:           aiDispatcher,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}
func (a *ReminderServiceImpl) ProcessAndSendReminders(ctx context.Context, reminderId uint) error {

	reminder, err := a.CommonService.GetModelById(ctx, reminderId, &[]string{"Request", "Request.User", "Request.Models", "Request.Models.AiModel", "Request.User.Telegram"})
	if err != nil {
		return err
	}
	for _, model := range *reminder.Request.Models {
		processor, err := a.AiDispatcher.RequestProcessor(ctx, model.AiModel, reminder.RequestID)
		if err != nil {
			return err
		}
		logger.Debug("ReminderServiceImpl.ProcessAndSendReminders", "processor", processor)
		if isCorrect, ok := processor["IsCorrect"].(bool); ok && isCorrect {

			var message string
			for key, value := range processor {
				if key != "IsCorrect" {
					message += key + ": " + fmt.Sprintf("%v", value) + "\n"
				}
			}
			return a.NotificationDispatcher.Notify(ctx, &types.Notification{
				Subject:  reminder.Message,
				Message:  message,
				Channels: []string{"sms", "email", "telegram"},
				UserId:   reminder.Request.UserID,
				User:     reminder.Request.User,
			})
		}

	}

	return nil
}
