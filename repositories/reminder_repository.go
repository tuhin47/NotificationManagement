package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type ReminderRepositoryImpl struct {
	*SQLRepository[models.Reminder]
}

func NewReminderRepository(db *gorm.DB) domain.ReminderRepository {
	return &ReminderRepositoryImpl{
		SQLRepository: NewSQLRepository[models.Reminder](db),
	}
}
