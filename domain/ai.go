package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"github.com/labstack/echo/v4"
)

type AIModelService interface {
	CommonService[models.AIModel]
}
type AIService[T any] interface {
	CommonService[T]
	MakeAIRequest(m *models.AIModel, requestId uint) (interface{}, error)
}
type AIModelRepository interface {
	Repository[models.AIModel, uint]
}

type AIProcessService[T AIService[X], X any] interface {
	ProcessAIRequest(*types.MakeAIRequestPayload) (interface{}, error)
	CreateModel(models.AIModelInterface) error
	UpdateModel(id uint, model models.AIModelInterface) (interface{}, error)
	GetModelById(id uint) (interface{}, error)
	GetAllModels(limit, offset int) (interface{}, error)
}

type AIRequestController interface {
	MakeAIRequestHandler(c echo.Context) error
	CreateAIModel(c echo.Context) error
	GetAIModelByID(c echo.Context) error
	GetAllAIModels(c echo.Context) error
	UpdateAIModel(c echo.Context) error
	DeleteAIModel(c echo.Context) error
}
