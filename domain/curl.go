package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"context"
	"github.com/labstack/echo/v4"
)

type CurlService interface {
	CommonService[models.CurlRequest]
	ProcessCurlRequest(c context.Context, req *models.CurlRequest) (*types.CurlResponse, error)
}

type CurlRequestRepository interface {
	Repository[models.CurlRequest, uint]
}
type AdditionalFieldsRepository interface {
	Repository[models.AdditionalFields, uint]
}

type CurlController interface {
	CurlHandler(c echo.Context) error
	GetCurlRequestByID(c echo.Context) error
	UpdateCurlRequest(c echo.Context) error
	DeleteCurlRequest(c echo.Context) error
}
