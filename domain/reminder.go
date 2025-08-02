package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type ReminderService interface {
	CommonService[models.Reminder]
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
