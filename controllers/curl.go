package controllers

import (
	"NotificationManagement/controllers/helper"
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CurlControllerImpl struct {
	CurlService domain.CurlService
	UserService domain.UserService
}

func NewCurlController(curlService domain.CurlService, userService domain.UserService) domain.CurlController {
	return &CurlControllerImpl{
		CurlService: curlService,
		UserService: userService,
	}
}

func (cc *CurlControllerImpl) CurlHandler(c echo.Context) error {
	var req types.CurlRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	if model.UserID == 0 {
		model.UserID = helper.GetUserId(c)
	}
	err = cc.CurlService.CreateModel(c.Request().Context(), model)
	if err != nil {
		return err
	}
	resp, err := cc.CurlService.ProcessCurlRequest(nil, model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"model":    model,
		"response": resp,
	})
}

func (cc *CurlControllerImpl) GetCurlRequestByID(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	curlRequest, err := cc.CurlService.GetModelById(c.Request().Context(), id, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, curlRequest)
}

func (cc *CurlControllerImpl) UpdateCurlRequest(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.CurlRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	if model.UserID == 0 {
		model.UserID = helper.GetUserId(c)
	}

	model, err = cc.CurlService.UpdateModel(c.Request().Context(), id, model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model)
}

func (cc *CurlControllerImpl) DeleteCurlRequest(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = cc.CurlService.DeleteModel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "CurlRequest deleted successfully"})
}
