package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"gorm.io/gorm"
)

type CurlRequestRepositoryImpl struct {
	domain.Repository[models.CurlRequest, uint]
}

func NewCurlRequestRepository(db *gorm.DB) domain.CurlRequestRepository {
	return &CurlRequestRepositoryImpl{
		Repository: NewSQLRepository[models.CurlRequest](db),
	}
}

type AdditionalFieldsRepositoryImpl struct {
	domain.Repository[models.AdditionalFields, uint]
}

func NewAdditionalFieldsRepository(db *gorm.DB) domain.AdditionalFieldsRepository {
	return &AdditionalFieldsRepositoryImpl{
		Repository: NewSQLRepository[models.AdditionalFields](db),
	}
}
