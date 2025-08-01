package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"errors"
	"fmt"
)

type AIServiceManagerImpl struct {
	deepseekService domain.DeepseekService
	geminiService   domain.GeminiService
	Repo            domain.AIModelRepository
}

func NewAIServiceManager(deepseekService domain.DeepseekService, geminiService domain.GeminiService, repo domain.AIModelRepository) domain.AIServiceManager {
	return &AIServiceManagerImpl{
		deepseekService: deepseekService,
		geminiService:   geminiService,
		Repo:            repo,
	}
}

func (f *AIServiceManagerImpl) ProcessAIRequest(req types.MakeAIRequestPayload) (interface{}, error) {
	model, err := f.GetModelByID(req.ModelID)
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

func (f *AIServiceManagerImpl) GetModelByID(id uint) (*models.AIModel, error) {
	return f.Repo.GetByID(context.Background(), id, nil)
}

func (f *AIServiceManagerImpl) GetService(modelType string) (interface{}, error) {
	switch modelType {
	case "deepseek", "local":
		return f.deepseekService, nil
	case "gemini":
		return f.geminiService, nil
	case "openai":
		// TODO: Implement OpenAI service when ready
		return nil, errors.New("OpenAI service not implemented yet")
	default:
		return nil, fmt.Errorf("unknown model type: %s", modelType)
	}
}
