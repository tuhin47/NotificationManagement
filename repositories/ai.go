package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type AIModelRepositoryImpl struct {
	domain.Repository[models.AIModel, uint]
}

func NewAIModelRepository(db *gorm.DB) domain.AIModelRepository {
	return &AIModelRepositoryImpl{
		Repository: NewSQLRepository[models.AIModel](db),
	}
}
