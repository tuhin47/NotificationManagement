package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DeepseekModelControllerImpl struct {
	Service domain.DeepseekModelService
}

func NewDeepseekModelController(service domain.DeepseekModelService) domain.DeepseekModelController {
	return &DeepseekModelControllerImpl{Service: service}
}

func (dc *DeepseekModelControllerImpl) CreateDeepseekModel(c echo.Context) error {
	var req types.DeepseekModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model := req.ToModel()
	err := dc.Service.CreateDeepseekModel(model)
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(model)
	return c.JSON(http.StatusCreated, response)
}

func (dc *DeepseekModelControllerImpl) GetDeepseekModelByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	model, err := dc.Service.GetDeepseekModelByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(model)
	return c.JSON(http.StatusOK, response)
}

func (dc *DeepseekModelControllerImpl) GetAllDeepseekModels(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	models, err := dc.Service.GetAllDeepseekModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []*types.DeepseekModelResponse
	for _, model := range models {
		responses = append(responses, types.FromDeepseekModel(&model))
	}

	return c.JSON(http.StatusOK, responses)
}

func (dc *DeepseekModelControllerImpl) UpdateDeepseekModel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	var req types.DeepseekModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model := req.ToModel()

	err = dc.Service.UpdateDeepseekModel(uint(id), model)
	if err != nil {
		return err
	}

	updatedModel, err := dc.Service.GetDeepseekModelByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(updatedModel)
	return c.JSON(http.StatusOK, response)
}

func (dc *DeepseekModelControllerImpl) DeleteDeepseekModel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	err = dc.Service.DeleteDeepseekModel(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "DeepseekModel deleted successfully"})
}
