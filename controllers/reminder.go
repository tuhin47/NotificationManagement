package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReminderControllerImpl struct {
	reminderService domain.ReminderService
	asynqService    domain.AsynqService
}

func NewReminderController(service domain.ReminderService, asynqService domain.AsynqService) domain.ReminderController {
	return &ReminderControllerImpl{reminderService: service, asynqService: asynqService}
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
	ctx := c.Request().Context()
	err = rc.reminderService.CreateModel(ctx, reminder)
	if err != nil {
		return err
	}
	reminder.TaskID, err = rc.asynqService.CreateReminderTask(ctx, reminder)
	if err != nil {
		return err
	}
	model, err := rc.reminderService.UpdateModel(ctx, reminder.ID, reminder)

	if err != nil {
		return err
	}

	response := types.FromReminderModel(model)
	return c.JSON(http.StatusCreated, response)
}

func (rc *ReminderControllerImpl) GetReminderByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	reminder, err := rc.reminderService.GetModelById(c.Request().Context(), id, nil)
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusOK, response)
}

func (rc *ReminderControllerImpl) GetAllReminders(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	reminders, err := rc.reminderService.GetAllModels(c.Request().Context(), limit, offset)
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
	reminder, err = rc.reminderService.UpdateModel(c.Request().Context(), id, reminder)
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

	err = rc.reminderService.DeleteModel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reminder deleted successfully"})
}
