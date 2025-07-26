package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"NotificationManagement/utils/throw"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

	reminder := req.ToModel()
	err := rc.Service.CreateReminder(reminder)
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusCreated, response)
}

func (rc *ReminderControllerImpl) GetReminderByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	reminder, err := rc.Service.GetReminderByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromReminderModel(reminder)
	return c.JSON(http.StatusOK, response)
}

func (rc *ReminderControllerImpl) GetAllReminders(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	reminders, err := rc.Service.GetAllReminders(limit, offset)
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	var req types.ReminderRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	reminder := req.ToModel()
	err = rc.Service.UpdateReminder(uint(id), reminder)
	if err != nil {
		return err
	}

	// Get the updated record
	updatedReminder, err := rc.Service.GetReminderByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromReminderModel(updatedReminder)
	return c.JSON(http.StatusOK, response)
}

func (rc *ReminderControllerImpl) DeleteReminder(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	err = rc.Service.DeleteReminder(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reminder deleted successfully"})
}
