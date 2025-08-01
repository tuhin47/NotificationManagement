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
	CommonService[T]
	MakeAIRequest(m *models.AIModel, requestId uint) (*Y, error)
}

type AIModelRepository interface {
	Repository[models.AIModel, uint]
}

type AIServiceManager interface {
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
