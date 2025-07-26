package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"github.com/labstack/echo/v4"
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

type CurlController interface {
	CurlHandler(c echo.Context) error
	GetCurlRequestByID(c echo.Context) error
	UpdateCurlRequest(c echo.Context) error
	DeleteCurlRequest(c echo.Context) error
}
