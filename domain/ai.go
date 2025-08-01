package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"github.com/labstack/echo/v4"
)

type AIModelType interface {
	models.DeepseekModel | models.GeminiModel
}

type AIService[T AIModelType, Y any] interface {
	MakeAIRequest(m *models.AIModel, requestId uint) (*Y, error)
	GetModelByID(id uint) (*T, error)
	CreateModel(model *T) error
	GetAllAIModels(limit, offset int) ([]T, error)
	UpdateAIModel(id uint, model *T) error
	DeleteAIModel(id uint) error
}

type AIModelRepository interface {
	Repository[models.AIModel, uint]
}

type AIServiceManager interface {
	GetService(modelType string) (interface{}, error)
	GetModelByID(id uint) (*models.AIModel, error)
	ProcessAIRequest(types.MakeAIRequestPayload) (interface{}, error)
}

type AIModelController interface {
	CreateAIModel(c echo.Context) error
	GetAIModelByID(c echo.Context) error
	GetAllAIModels(c echo.Context) error
	UpdateAIModel(c echo.Context) error
	DeleteAIModel(c echo.Context) error
}

type AIRequestController interface {
	MakeAIRequestHandler(c echo.Context) error
}
