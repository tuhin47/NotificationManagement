package controllers

import (
	"net/http"
	"strconv"

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

func (cc *CurlController) GetCurlRequestByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	curlRequest, err := cc.Service.GetCurlRequestByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "CurlRequest not found"})
	}

	return c.JSON(http.StatusOK, curlRequest)
}
