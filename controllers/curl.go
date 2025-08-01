package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CurlControllerImpl struct {
	Service domain.CurlService
}

func NewCurlController(service domain.CurlService) domain.CurlController {
	return &CurlControllerImpl{Service: service}
}

func (cc *CurlControllerImpl) CurlHandler(c echo.Context) error {
	var req types.CurlRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	err = cc.Service.CreateModel(model)
	if err != nil {
		return err
	}
	resp, err := cc.Service.ExecuteCurl(model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, *resp)
}

func (cc *CurlControllerImpl) GetCurlRequestByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	curlRequest, err := cc.Service.GetModelByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, curlRequest)
}

func (cc *CurlControllerImpl) UpdateCurlRequest(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.CurlRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model, err := cc.Service.UpdateCurlRequest(id, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model)
}

func (cc *CurlControllerImpl) DeleteCurlRequest(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = cc.Service.DeleteModel(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "CurlRequest deleted successfully"})
}
