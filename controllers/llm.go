package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LLMControllerImpl struct {
	Service domain.LLMService
}

func NewLLMController(service domain.LLMService) domain.LLMController {
	return &LLMControllerImpl{Service: service}
}

func (lc *LLMControllerImpl) CreateLLM(c echo.Context) error {
	var req types.LLMRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm := req.ToModel()
	err := lc.Service.CreateModel(llm)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusCreated, response)
}

func (lc *LLMControllerImpl) GetLLMByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	llm, err := lc.Service.GetModelByID(id)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) GetAllLLMs(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	llms, err := lc.Service.GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []*types.LLMResponse
	for _, llm := range llms {
		responses = append(responses, types.FromLLMModel(&llm))
	}

	return c.JSON(http.StatusOK, responses)
}

func (lc *LLMControllerImpl) UpdateLLM(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.LLMRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm := req.ToModel()
	err = lc.Service.UpdateModel(id, llm)
	if err != nil {
		return err
	}

	// Get the updated record
	updatedLLM, err := lc.Service.GetModelByID(id)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(updatedLLM)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) DeleteLLM(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	err = lc.Service.DeleteModel(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "LLM deleted successfully"})
}
