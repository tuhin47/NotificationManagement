package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"github.com/labstack/echo/v4"
)

type BaseAIProcessImpl[T domain.AIService[X], X any] struct {
	domain.CommonService[X]
	Service        T
	AIModelService domain.AIModelService
}

func NewAIServiceManager[T domain.AIService[X], X any](aiService domain.AIModelService, service T) domain.AIProcessService[T, X] {
	return &BaseAIProcessImpl[T, X]{
		AIModelService: aiService,
		Service:        service,
	}
}

func (f *BaseAIProcessImpl[T, X]) CreateModel(c echo.Context, model models.AIModelInterface) error {
	x := any(model).(*X)
	return f.Service.CreateModel(c, x)
}
func (f *BaseAIProcessImpl[T, X]) GetModelById(c echo.Context, id uint) (interface{}, error) {
	return f.Service.GetModelById(c, id)
}
func (f *BaseAIProcessImpl[T, X]) GetAllModels(c echo.Context, limit int, offset int) (interface{}, error) {
	return f.Service.GetAllModels(c, limit, offset)
}

func (f *BaseAIProcessImpl[T, X]) UpdateModel(c echo.Context, id uint, model models.AIModelInterface) (interface{}, error) {
	x := any(model).(*X)
	return f.Service.UpdateModel(c, id, x)
}

func (f *BaseAIProcessImpl[T, X]) MakeAIRequest(c echo.Context, req *types.MakeAIRequestPayload) (interface{}, error) {
	model, err := f.AIModelService.GetModelById(c, req.ModelID)
	if err != nil {
		return nil, err
	}
	resp, err := f.Service.MakeAIRequest(c, model, req.CurlRequestID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type DeepseekProcessServiceImpl struct {
	domain.AIProcessService[domain.AIService[models.DeepseekModel], models.DeepseekModel]
}

func NewDeepseekServiceManager(aiService domain.AIModelService, service domain.DeepseekService) DeepseekProcessServiceImpl {
	return DeepseekProcessServiceImpl{
		AIProcessService: NewAIServiceManager[domain.AIService[models.DeepseekModel]](aiService, service),
	}
}

type GeminiProcessServiceImpl struct {
	domain.AIProcessService[domain.AIService[models.GeminiModel], models.GeminiModel]
}

func NewGeminiServiceManager(aiService domain.AIModelService, service domain.GeminiService) GeminiProcessServiceImpl {
	return GeminiProcessServiceImpl{
		AIProcessService: NewAIServiceManager[domain.AIService[models.GeminiModel]](aiService, service),
	}
}
