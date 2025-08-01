package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AIModelControllerImpl struct {
	Service domain.DeepseekService
}

func NewAIModelController(service domain.DeepseekService) domain.AIModelController {
	return &AIModelControllerImpl{Service: service}
}

func (dc *AIModelControllerImpl) CreateAIModel(c echo.Context) error {
	var req types.AIModelRequest
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
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	model, err := dc.Service.GetModelByID(id)
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(model)
	return c.JSON(http.StatusOK, response)
}

func (dc *AIModelControllerImpl) GetAllAIModels(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	models, err := dc.Service.GetAllModels(limit, offset)
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
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.AIModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model := req.ToModel()

	err = dc.Service.UpdateModel(id, model)
	if err != nil {
		return err
	}

	updatedModel, err := dc.Service.GetModelByID(id)
	if err != nil {
		return err
	}

	response := types.FromDeepseekModel(updatedModel)
	return c.JSON(http.StatusOK, response)
}

func (dc *AIModelControllerImpl) DeleteAIModel(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = dc.Service.DeleteModel(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "DeepseekModel deleted successfully"})
}
