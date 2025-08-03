package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type GeminiRepositoryImpl struct {
	domain.Repository[models.GeminiModel, uint]
}

func NewGeminiRepository(db *gorm.DB) domain.GeminiModelRepository {
	return &GeminiRepositoryImpl{
		Repository: NewSQLRepository[models.GeminiModel](db),
	}
}
