package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReminderControllerImpl struct {
	Service domain.ReminderService
}

func NewReminderController(service domain.ReminderService) domain.ReminderController {
	return &ReminderControllerImpl{Service: service}
}

func (rc *ReminderControllerImpl) CreateReminder(c echo.Context) error {
	var req types.ReminderRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	reminder, err := req.ToModel()
	if err != nil {
		return err
	}
	err = rc.Service.CreateModel(reminder)
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusCreated, response)
}

func (rc *ReminderControllerImpl) GetReminderByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	reminder, err := rc.Service.GetModelById(id)
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusOK, response)
}

func (rc *ReminderControllerImpl) GetAllReminders(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	reminders, err := rc.Service.GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []*types.ReminderResponse
	for _, reminder := range reminders {
		responses = append(responses, types.FromReminderModel(&reminder))
	}

	return c.JSON(http.StatusOK, responses)
}

func (rc *ReminderControllerImpl) UpdateReminder(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.ReminderRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	reminder, err := req.ToModel()
	if err != nil {
		return err
	}
	reminder, err = rc.Service.UpdateModel(id, reminder)
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusOK, response)
}

func (rc *ReminderControllerImpl) DeleteReminder(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = rc.Service.DeleteModel(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reminder deleted successfully"})
}
