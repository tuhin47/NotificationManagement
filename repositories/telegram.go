package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"
	"gorm.io/gorm"
)

type TelegramRepositoryImpl struct {
	domain.Repository[models.Telegram, uint]
}

func NewTelegramRepository(db *gorm.DB) domain.TelegramRepository {
	return &TelegramRepositoryImpl{
		Repository: NewSQLRepository[models.Telegram](db),
	}
}

func (r *TelegramRepositoryImpl) GetByOTP(ctx context.Context, otp string) (*models.Telegram, error) {

	var telegram models.Telegram
	err := r.Repository.GetDB(ctx).Where("otp = ?", otp).First(&telegram).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &telegram, nil
}
