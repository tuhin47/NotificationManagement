package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type ReminderServiceImpl struct {
	domain.CommonService[models.Reminder]
	Repo domain.ReminderRepository
}

func NewReminderService(repo domain.ReminderRepository) domain.ReminderService {
	service := &ReminderServiceImpl{
		Repo: repo,
	}
	service.CommonService = NewCommonService[models.Reminder](repo, service)
	return service
}
