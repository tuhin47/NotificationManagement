package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type DeepseekModelRepositoryImpl struct {
	domain.Repository[models.DeepseekModel, uint]
}

func NewDeepseekModelRepository(db *gorm.DB) domain.DeepseekModelRepository {
	return &DeepseekModelRepositoryImpl{
		Repository: NewSQLRepository[models.DeepseekModel](db),
	}
}
