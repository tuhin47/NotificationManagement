package controllers

import (
	"NotificationManagement/controllers/helper"
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIRequestControllerImpl struct {
	domain.AIModelService
	domain.AiDispatcher
}

func NewAIRequestController(aiModelService domain.AIModelService, service domain.AiDispatcher) domain.AIRequestController {
	return &AIRequestControllerImpl{AIModelService: aiModelService, AiDispatcher: service}
}

func (a *AIRequestControllerImpl) MakeAIRequestHandler(c echo.Context) error {
	var req types.MakeAIRequestPayload
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}
	context := c.Request().Context()
	model, err := a.GetModelById(context, req.ModelID, nil)
	if err != nil {
		return err
	}
	aiResponse, err := a.RequestProcessor(context, model, req.CurlRequestID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, aiResponse)
}

func (a *AIRequestControllerImpl) CreateAIModel(c echo.Context) error {

	var req types.AIModelRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	err = a.ProcessCreateModel(c.Request().Context(), model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, model)

}

func (a *AIRequestControllerImpl) UpdateAIModel(c echo.Context) error {
	var req types.AIModelRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	req.ID = id
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	updateModel, err := a.ProcessUpdateModel(c.Request().Context(), model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, updateModel)
}

func (a *AIRequestControllerImpl) GetAIModelByID(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	modelById, err := a.ProcessModelById(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelById)
}

func (a *AIRequestControllerImpl) GetAllAIModels(c echo.Context) error {
	//limit, offset := helper.ParseLimitAndOffset(c)
	var responses []any

	responses = a.ProcessAllAIModels(c.Request().Context())

	return c.JSON(http.StatusOK, responses)
}

func (a *AIRequestControllerImpl) DeleteAIModel(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	return a.DeleteModel(c.Request().Context(), id)
}
