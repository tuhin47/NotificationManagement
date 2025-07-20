package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
)

type CurlService interface {
	ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error)
	GetCurlRequestByID(id uint) (*models.CurlRequest, error)
	UpdateCurlRequest(id uint, req types.CurlRequest) (*models.CurlRequest, error)
	DeleteCurlRequest(id uint) error
}

type CurlRequestRepository interface {
	Repository[models.CurlRequest, uint]
}
