package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationController struct {
	NotificationService domain.NotificationService
}

func NewNotificationController(notificationService domain.NotificationService) *NotificationController {
	return &NotificationController{NotificationService: notificationService}
}

type NotifyRequest struct {
	To       string   `json:"to"`
	Subject  string   `json:"subject"`
	Message  string   `json:"message"`
	Channels []string `json:"channels"`
}

func (h *NotificationController) Notify(c echo.Context) error {
	var req NotifyRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	err := h.NotificationService.Send(&types.Notification{
		To:       req.To,
		Subject:  req.Subject,
		Message:  req.Message,
		Channels: req.Channels,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Notifications sent"})
}
