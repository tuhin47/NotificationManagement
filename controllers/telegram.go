package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/middleware"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TelegramControllerImpl struct {
	telegramService domain.TelegramNotifier
}

func NewTelegramController(service domain.TelegramNotifier) domain.TelegramController {
	return &TelegramControllerImpl{telegramService: service}
}

func (t TelegramControllerImpl) VerifyOtp(c echo.Context) error {
	ccx, _ := c.(*middleware.CustomContext)

	var req types.VerifyOtpRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}
	userID := ccx.UserID

	// TODO- Consider some expiration in future
	telegramModel, err := t.telegramService.VerifyOTP(c.Request().Context(), req.OTP, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "OTP verified successfully",
		"telegram": telegramModel,
	})
}
