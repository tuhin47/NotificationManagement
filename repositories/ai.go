package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type AIModelRepositoryImpl struct {
	*SQLRepository[models.AIModel]
}

func NewAIModelRepository(db *gorm.DB) domain.AIModelRepository {
	return &AIModelRepositoryImpl{
		SQLRepository: NewSQLRepository[models.AIModel](db),
	}
}
