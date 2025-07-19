package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
)

type CurlService interface {
	ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error)
}

type CurlRequestRepository interface {
	Repository[models.CurlRequest, uint]
}
