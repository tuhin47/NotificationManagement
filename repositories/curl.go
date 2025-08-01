package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"gorm.io/gorm"
)

type CurlRequestRepositoryImpl struct {
	*SQLRepository[models.CurlRequest]
}

func NewCurlRequestRepository(db *gorm.DB) domain.CurlRequestRepository {
	return &CurlRequestRepositoryImpl{
		SQLRepository: NewSQLRepository[models.CurlRequest](db),
	}
}

type AdditionalFieldsRepositoryImpl struct {
	*SQLRepository[models.AdditionalFields]
}

func NewAdditionalFieldsRepository(db *gorm.DB) domain.AdditionalFieldsRepository {
	return &AdditionalFieldsRepositoryImpl{
		SQLRepository: NewSQLRepository[models.AdditionalFields](db),
	}
}
