package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/services/notifier"
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

func (h *NotificationController) NotifyAll(c echo.Context) error {
	var req NotifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	err := h.NotificationService.Send(notifier.Notification{
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
