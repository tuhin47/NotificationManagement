package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type GeminiRepositoryImpl struct {
	*SQLRepository[models.GeminiModel]
}

func NewGeminiRepository(db *gorm.DB) domain.GeminiModelRepository {
	return &GeminiRepositoryImpl{
		SQLRepository: NewSQLRepository[models.GeminiModel](db),
	}
}
