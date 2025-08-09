package controllers

import (
	"NotificationManagement/controllers/helper"
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationController struct {
	domain.NotificationDispatcher
}

func NewNotificationController(notificationService domain.NotificationDispatcher) *NotificationController {
	return &NotificationController{NotificationDispatcher: notificationService}
}

func (h *NotificationController) Notify(c echo.Context) error {
	var req types.NotifyRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}

	err := h.NotificationDispatcher.Notify(c.Request().Context(), &types.Notification{
		UserId:   req.UserId,
		Subject:  req.Subject,
		Message:  req.Message,
		Channels: req.Channels,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Notifications sent"})
}
