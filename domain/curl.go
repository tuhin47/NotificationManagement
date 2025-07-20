package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
)

type CurlService interface {
	ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error)
	GetCurlRequestByID(id uint) (*models.CurlRequest, error)
}

type CurlRequestRepository interface {
	Repository[models.CurlRequest, uint]
}
