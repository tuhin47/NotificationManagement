package controllers

import (
	"NotificationManagement/utils/errutil"
	"net/http"
	"strconv"

	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils/throw"

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
		return throw.AppError(errutil.ErrInvalidRequestBody, err)
	}
	resp, err := cc.Service.ExecuteCurl(req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (cc *CurlController) GetCurlRequestByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	curlRequest, err := cc.Service.GetCurlRequestByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, curlRequest)
}

func (cc *CurlController) UpdateCurlRequest(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	var req types.CurlRequest
	if err := c.Bind(&req); err != nil {
		return throw.AppError(errutil.ErrInvalidRequestBody, err)
	}

	updatedRequest, err := cc.Service.UpdateCurlRequest(uint(id), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedRequest)
}

func (cc *CurlController) DeleteCurlRequest(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return throw.AppError(errutil.ErrInvalidIdParam, err)
	}

	err = cc.Service.DeleteCurlRequest(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "CurlRequest deleted successfully"})
}
