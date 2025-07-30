package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIControllerImpl struct {
	OllamaService   domain.OllamaService
	DeepseekService domain.DeepseekModelService
}

func NewAIController(ollamaService domain.OllamaService, deepseekService domain.DeepseekModelService) domain.AIController {
	return &AIControllerImpl{
		OllamaService:   ollamaService,
		DeepseekService: deepseekService,
	}
}

func (ac *AIControllerImpl) MakeAIRequestHandler(c echo.Context) error {
	var req types.MakeAIRequestPayload
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	deepseekModel, err := ac.DeepseekService.GetDeepseekModelByID(req.ModelID)
	if err != nil {
		return err
	}

	ollamaResponse, err := ac.OllamaService.MakeAIRequest(deepseekModel, req.CurlRequestID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, *ollamaResponse)
}
