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
	err := lc.Service.CreateLLM(llm)
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusCreated, response)
}

func (lc *LLMControllerImpl) GetLLMByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	llm, err := lc.Service.GetLLMByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromLLMModel(llm)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) GetAllLLMs(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10 // default limit
	offset := 0 // default offset

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

	llms, err := lc.Service.GetAllLLMs(limit, offset)
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	var req types.LLMRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	llm := req.ToModel()
	err = lc.Service.UpdateLLM(uint(id), llm)
	if err != nil {
		return err
	}

	// Get the updated record
	updatedLLM, err := lc.Service.GetLLMByID(uint(id))
	if err != nil {
		return err
	}

	response := types.FromLLMModel(updatedLLM)
	return c.JSON(http.StatusOK, response)
}

func (lc *LLMControllerImpl) DeleteLLM(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}

	err = lc.Service.DeleteLLM(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "LLM deleted successfully"})
}
