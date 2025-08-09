package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
)

type AiDispatcherImpl struct {
	services *[]domain.DispatchableAIService
}

func NewAIDispatcher(geminiService domain.GeminiService, deepseekService domain.DeepseekService) domain.AiDispatcher {
	return &AiDispatcherImpl{
		services: &[]domain.DispatchableAIService{
			geminiService,
			deepseekService,
		},
	}
}

func (a AiDispatcherImpl) RequestProcessor(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error) {
	for _, service := range *a.services {
		if service.GetModelType() == m.Type {
			return service.GetAIJsonResponse(c, m, requestId)
		}
	}
	return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}
