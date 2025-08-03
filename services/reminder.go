package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type ReminderServiceImpl struct {
	domain.CommonService[models.Reminder]
}

func NewReminderService(repo domain.ReminderRepository) domain.ReminderService {
	service := &ReminderServiceImpl{}
	service.CommonService = NewCommonService(repo, service)
	return service
}
