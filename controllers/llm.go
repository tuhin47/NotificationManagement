package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
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
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm, err := req.ToModel()
	if err != nil {
		return err
	}
	err = lc.Service.CreateModel(c, llm)
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

	llm, err := lc.Service.GetModelById(c, id)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) GetAllLLMs(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	llms, err := lc.Service.GetAllModels(c, limit, offset)
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

	llm, err := req.ToModel()
	if err != nil {
		return err
	}
	updatedLLM, err := lc.Service.UpdateModel(c, id, llm)
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

	err = lc.Service.DeleteModel(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "LLM deleted successfully"})
}
