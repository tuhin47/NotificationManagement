package controllers

import (
	"net/http"

	"NotificationManagement/domain"
	"NotificationManagement/types"

	"github.com/labstack/echo/v4"
)

type CurlController struct {
	Service domain.CurlService
}

func NewCurlController(service domain.CurlService) *CurlController {
	return &CurlController{Service: service}
}

// Handler for POST /api/curl
func (cc *CurlController) CurlHandler(c echo.Context) error {
	var req types.CurlRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.CurlResponse{ErrMessage: "Invalid request payload"})
	}

	resp, err := cc.Service.ExecuteCurl(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, resp)
	}
	return c.JSON(http.StatusOK, resp)
}
