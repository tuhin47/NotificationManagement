package notifier

import (
	"NotificationManagement/config"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"fmt"
	"strings"
)

type TelegramNotifierImpl struct {
	//BotAPI       *tgbotapi.BotAPI
	telegramRepo domain.TelegramRepository
}

func NewTelegramNotifier(repo domain.TelegramRepository) domain.TelegramNotifier {
	return &TelegramNotifierImpl{
		telegramRepo: repo,
	}
}

func (t *TelegramNotifierImpl) Send(ctx context.Context, notification *types.Notification) error {
	/*
		TODO : have to call the worker
		chatID := (*notification.User.Telegram)[0].ChatID
		msg := tgbotapi.NewMessage(chatID, notification.Message)
		_, err := t.BotAPI.Send(msg)
		if err != nil {
			logger.Error("Failed to send Telegram message", "error", err)
			return err
		}
		logger.Debug(fmt.Sprintf("[Telegram] To: %s, Message: %s", chatID, notification.Message), notification)
	*/
	logger.Debug(notification.Message)
	return nil
}

func (t *TelegramNotifierImpl) Type() string {
	return "telegram"
}

func (t *TelegramNotifierImpl) IsActive() bool {
	return *config.Telegram().Enabled
}

func (t *TelegramNotifierImpl) VerifyOTP(ctx context.Context, otp string, userID uint) (*models.Telegram, error) {
	telegramModel, err := t.telegramRepo.GetByOTP(ctx, otp)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrRecordNotFound, fmt.Errorf("telegram chat ID not found"))
	}

	if strings.EqualFold(telegramModel.Otp, otp) {
		telegramModel.UserID = &userID
		err = t.telegramRepo.Update(ctx, telegramModel)
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrDatabaseQuery, fmt.Errorf("failed to update telegram user ID: %w", err))
		}
		return telegramModel, nil
	}
	return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, fmt.Errorf("invalid OTP"))
}
