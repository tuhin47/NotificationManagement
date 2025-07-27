package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"github.com/labstack/echo/v4"
)

type AIService[T any, Y any] interface {
	MakeAIRequest(mod T, requestId uint) (Y, error)
}

type OllamaService interface {
	AIService[*models.DeepseekModel, *types.OllamaResponse]
	PullModel(model models.DeepseekModel) error
}

type AIController interface {
	MakeAIRequestHandler(c echo.Context) error
}
