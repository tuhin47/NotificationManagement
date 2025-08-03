package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type ReminderRepositoryImpl struct {
	domain.Repository[models.Reminder, uint]
}

func NewReminderRepository(db *gorm.DB) domain.ReminderRepository {
	return &ReminderRepositoryImpl{
		Repository: NewSQLRepository[models.Reminder](db),
	}
}
