package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
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
func (f *BaseAIProcessImpl[T, X]) CreateModel(model models.AIModelInterface) error {
	x := any(model).(*X)
	return f.Service.CreateModel(x)
}
func (f *BaseAIProcessImpl[T, X]) GetModelById(id uint) (interface{}, error) {
	return f.Service.GetModelById(id)
}
func (f *BaseAIProcessImpl[T, X]) GetAllModels(limit int, offset int) (interface{}, error) {
	return f.Service.GetAllModels(limit, offset)
}

func (f *BaseAIProcessImpl[T, X]) UpdateModel(id uint, model models.AIModelInterface) (interface{}, error) {
	x := any(model).(*X)
	return f.Service.UpdateModel(id, x)
}

func (f *BaseAIProcessImpl[T, X]) MakeAIRequest(req *types.MakeAIRequestPayload) (interface{}, error) {
	model, err := f.AIModelService.GetModelById(req.ModelID)
	if err != nil {
		return nil, err
	}
	resp, err := f.Service.MakeAIRequest(model, req.CurlRequestID)
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
