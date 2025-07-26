package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

	"gorm.io/gorm"
)

type DeepseekModelRepositoryImpl struct {
	*SQLRepository[models.DeepseekModel]
}

func NewDeepseekModelRepository(db *gorm.DB) domain.DeepseekModelRepository {
	return &DeepseekModelRepositoryImpl{
		SQLRepository: NewSQLRepository[models.DeepseekModel](db),
	}
}
