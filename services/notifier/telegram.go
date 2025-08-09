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
	domain.TelegramAPI
	telegramRepo domain.TelegramRepository
}

func NewTelegramNotifier(repo domain.TelegramRepository, api domain.TelegramAPI) domain.TelegramNotifier {
	return &TelegramNotifierImpl{
		telegramRepo: repo,
		TelegramAPI:  api,
	}
}

func (t *TelegramNotifierImpl) Send(ctx context.Context, notification *types.Notification) error {

	chatID := (*notification.User.Telegram)[0].ChatID
	//TODO : have create an call  worker
	if *config.Telegram().Enabled {
		t.TelegramAPI.SendMessage(chatID, notification.Message, nil)
		logger.Info(fmt.Sprintf("[Telegram] To: %s, Message: %s", chatID, notification.Message), notification)
	}

	logger.Info(fmt.Sprintf("[Telegram] To: %s, Message: %s", chatID, notification.Message), notification)

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
