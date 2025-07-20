package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"gorm.io/gorm"
)

type GormCurlRequestRepository struct {
	*SQLRepository[models.CurlRequest]
}

func NewCurlRequestRepository(db *gorm.DB) domain.CurlRequestRepository {
	return &GormCurlRequestRepository{
		SQLRepository: NewSQLRepository[models.CurlRequest](db),
	}
}
