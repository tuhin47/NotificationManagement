package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"context"

	"github.com/labstack/echo/v4"
)

type AIModelService interface {
	CommonService[models.AIModel]
}

type AIService[T any] interface {
	CommonService[T]
	MakeAIRequest(c context.Context, m *models.AIModel, requestId uint) (interface{}, error)
	GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error)
	GetModelType() string
}
type AIModelRepository interface {
	Repository[models.AIModel, uint]
}

type AiDispatcher interface {
	RequestProcessor(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error)
}

type DispatchableAIService interface {
	GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error)
	GetModelType() string
}

type AIProcessService[T AIService[X], X any] interface {
	MakeAIRequest(c context.Context, req *types.MakeAIRequestPayload) (interface{}, error)
	CreateModel(c context.Context, model models.AIModelInterface) error
	UpdateModel(c context.Context, id uint, model models.AIModelInterface) (interface{}, error)
	GetModelById(c context.Context, id uint) (interface{}, error)
	GetAllModels(c context.Context, limit, offset int) (interface{}, error)
}

type AIRequestController interface {
	MakeAIRequestHandler(c echo.Context) error
	CreateAIModel(c echo.Context) error
	GetAIModelByID(c echo.Context) error
	GetAllAIModels(c echo.Context) error
	UpdateAIModel(c echo.Context) error
	DeleteAIModel(c echo.Context) error
}
