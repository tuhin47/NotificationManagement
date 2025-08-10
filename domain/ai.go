package domain

import (
	"NotificationManagement/models"
	"context"

	"github.com/labstack/echo/v4"
)

type AIModelService interface {
	CommonService[models.AIModel]
}

type AIService[T any] interface {
	DispatchableAIService
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
	ProcessCreateModel(ctx context.Context, model models.AIModelInterface) error
	ProcessModelById(ctx context.Context, id uint) (any, error)
	ProcessAllAIModels(ctx context.Context) []any
	ProcessUpdateModel(ctx context.Context, model models.AIModelInterface) (any, error)
}

type DispatchableAIService interface {
	GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error)
	GetModelType() string
	CreateAIModel(c context.Context, model any) error
	GetAIModelById(ctx context.Context, id uint) (any, error)
	GetAllAIModels(ctx context.Context) ([]any, error)
	UpdateAIModel(c context.Context, model any) (any, error)
}

type AIRequestController interface {
	MakeAIRequestHandler(c echo.Context) error
	CreateAIModel(c echo.Context) error
	GetAIModelByID(c echo.Context) error
	GetAllAIModels(c echo.Context) error
	UpdateAIModel(c echo.Context) error
	DeleteAIModel(c echo.Context) error
}
