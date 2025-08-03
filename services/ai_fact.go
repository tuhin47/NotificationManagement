package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"encoding/json"
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
		resp, err := f.geminiService.MakeAIRequest(model, req.CurlRequestID)
		if err != nil {
			return nil, err
		}
		var result interface{}
		if err := json.Unmarshal([]byte(resp.Text()), &result); err != nil {
			return nil, err
		}
		return result, nil

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
