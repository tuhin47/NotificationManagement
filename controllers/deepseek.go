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

type AIModelControllerImpl struct {
	Service domain.DeepseekService
}

func NewAIModelController(service domain.DeepseekService) domain.AIModelController {
	return &AIModelControllerImpl{Service: service}
}

func (dc *AIModelControllerImpl) CreateAIModel(c echo.Context) error {
	var req types.DeepseekModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model := req.ToModel()
	err := dc.Service.CreateModel(model)
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(model)
	return c.JSON(http.StatusCreated, response)
}

func (dc *AIModelControllerImpl) GetAIModelByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	model, err := dc.Service.GetModelByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(model)
	return c.JSON(http.StatusOK, response)
}

func (dc *AIModelControllerImpl) GetAllAIModels(c echo.Context) error {
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

	models, err := dc.Service.GetAllAIModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []*types.DeepseekModelResponse
	for _, model := range models {
		responses = append(responses, types.FromDeepseekModel(&model))
	}

	return c.JSON(http.StatusOK, responses)
}

func (dc *AIModelControllerImpl) UpdateAIModel(c echo.Context) error {
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

	err = dc.Service.UpdateAIModel(uint(id), model)
	if err != nil {
		return err
	}

	updatedModel, err := dc.Service.GetModelByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(updatedModel)
	return c.JSON(http.StatusOK, response)
}

func (dc *AIModelControllerImpl) DeleteAIModel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	err = dc.Service.DeleteAIModel(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "DeepseekModel deleted successfully"})
}
