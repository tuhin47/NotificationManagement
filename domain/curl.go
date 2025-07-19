package domain

import "NotificationManagement/types"

type CurlService interface {
	ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error)
}
