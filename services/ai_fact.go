package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
)

type AIServiceManagerImpl struct {
	deepseekService domain.DeepseekService
	geminiService   domain.GeminiService
	AIModelRepo     domain.AIModelRepository
}

func NewAIServiceManager(deepseekService domain.DeepseekService, geminiService domain.GeminiService, repo domain.AIModelRepository) domain.AIServiceManager {
	return &AIServiceManagerImpl{
		deepseekService: deepseekService,
		geminiService:   geminiService,
		AIModelRepo:     repo,
	}
}

func (f *AIServiceManagerImpl) ProcessAIRequest(req types.MakeAIRequestPayload) (interface{}, error) {
	model, err := f.AIModelRepo.GetByID(context.Background(), req.ModelID, nil)
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
