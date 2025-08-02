package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
)

type AIServiceManagerImpl struct {
	deepseekService domain.DeepseekService
	geminiService   domain.GeminiService
	AIModelService  domain.AIModelService
}

func NewAIServiceManager(deepseekService domain.DeepseekService, geminiService domain.GeminiService, aiService domain.AIModelService) domain.AIServiceManager {
	return &AIServiceManagerImpl{
		deepseekService: deepseekService,
		geminiService:   geminiService,
		AIModelService:  aiService,
	}
}

func (f *AIServiceManagerImpl) ProcessAIRequest(req types.MakeAIRequestPayload) (interface{}, error) {
	model, err := f.AIModelService.GetModelByID(req.ModelID)
	if err != nil {
		return nil, err
	}
	switch model.Type {
	case "deepseek", "local":
		return f.deepseekService.MakeAIRequest(model, req.CurlRequestID)

	case "gemini":
		return f.geminiService.MakeAIRequest(model, req.CurlRequestID)

	default:
		return nil, errutil.NewAppError(errutil.ErrServiceUnavailable, errutil.ErrServiceNotAvailable)
	}
}

func (f *AIServiceManagerImpl) GetService(modelType string) (interface{}, error) {
	switch modelType {
	case "deepseek", "local":
		return f.deepseekService, nil
	case "gemini":
		return f.geminiService, nil
	default:
		return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
	}
}
