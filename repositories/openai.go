package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type OpenAIModelRepositoryImpl struct {
	domain.Repository[models.OpenAIModel, uint]
}

func NewOpenAIModelRepository(db *gorm.DB) domain.OpenAIModelRepository {
	return &OpenAIModelRepositoryImpl{
		Repository: NewSQLRepository[models.OpenAIModel](db),
	}
}
