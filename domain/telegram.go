package domain

import (
	"NotificationManagement/models"
	"context"

	"github.com/labstack/echo/v4"
)

type TelegramAPI interface {
	Start()
	SendMessage(chatID int64, text string, markup interface{})
}

type TelegramNotifier interface {
	Notifier
	VerifyOTP(ctx context.Context, otp string, userID uint) (*models.Telegram, error)
}

type TelegramRepository interface {
	Repository[models.Telegram, uint]
	GetByOTP(ctx context.Context, otp string) (*models.Telegram, error)
	GetByChatId(ctx context.Context, id int64) (*models.Telegram, error)
}

type TelegramController interface {
	VerifyOtp(c echo.Context) error
}
