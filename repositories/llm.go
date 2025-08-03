package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type LLMRepositoryImpl struct {
	*SQLRepository[models.RequestAIModel]
}

func NewLLMRepository(db *gorm.DB) domain.LLMRepository {
	return &LLMRepositoryImpl{
		SQLRepository: NewSQLRepository[models.RequestAIModel](db),
	}
}
