package controllers

import (
	"NotificationManagement/controllers/helper"
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LLMControllerImpl struct {
	Service domain.LLMService
}

func NewLLMController(service domain.LLMService) domain.LLMController {
	return &LLMControllerImpl{Service: service}
}

func (lc *LLMControllerImpl) CreateLLM(c echo.Context) error {
	var req types.LLMRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm, err := req.ToModel()
	if err != nil {
		return err
	}
	err = lc.Service.CreateModel(c.Request().Context(), llm)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusCreated, response)
}

func (lc *LLMControllerImpl) GetLLMByID(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	llm, err := lc.Service.GetModelById(c.Request().Context(), id, nil)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) GetAllLLMs(c echo.Context) error {
	limit, offset := helper.ParseLimitAndOffset(c)

	models, err := lc.Service.GetAllModels(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}

	var responses []*types.LLMResponse
	for _, llm := range models {
		responses = append(responses, types.FromLLMModel(&llm))
	}

	return c.JSON(http.StatusOK, responses)
}

func (lc *LLMControllerImpl) UpdateLLM(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.LLMRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm, err := req.ToModel()
	if err != nil {
		return err
	}
	updatedLLM, err := lc.Service.UpdateModel(c.Request().Context(), id, llm)
	if err != nil {
		return err
	}
	response := types.FromLLMModel(updatedLLM)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) DeleteLLM(c echo.Context) error {
	id, err := helper.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = lc.Service.DeleteModel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "LLM deleted successfully"})
}
