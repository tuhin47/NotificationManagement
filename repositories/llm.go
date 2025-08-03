package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type LLMRepositoryImpl struct {
	domain.Repository[models.RequestAIModel, uint]
}

func NewLLMRepository(db *gorm.DB) domain.LLMRepository {
	return &LLMRepositoryImpl{
		Repository: NewSQLRepository[models.RequestAIModel](db),
	}
}
