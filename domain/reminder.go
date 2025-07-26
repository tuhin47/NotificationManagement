package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type ReminderService interface {
	CreateReminder(reminder *models.Reminder) error
	GetReminderByID(id uint) (*models.Reminder, error)
	GetAllReminders(limit, offset int) ([]models.Reminder, error)
	UpdateReminder(id uint, reminder *models.Reminder) error
	DeleteReminder(id uint) error
}

type ReminderRepository interface {
	Repository[models.Reminder, uint]
}

type ReminderController interface {
	CreateReminder(c echo.Context) error
	GetReminderByID(c echo.Context) error
	GetAllReminders(c echo.Context) error
	UpdateReminder(c echo.Context) error
	DeleteReminder(c echo.Context) error
}
