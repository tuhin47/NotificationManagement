package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
)

type AiDispatcherImpl struct {
	services *[]domain.DispatchableAIService
	aiModel  domain.AIModelService
}

func NewAIDispatcher(geminiService domain.GeminiService, deepseekService domain.DeepseekService, openaiService domain.OpenAIService, ai domain.AIModelService) domain.AiDispatcher {
	return &AiDispatcherImpl{
		services: &[]domain.DispatchableAIService{
			geminiService,
			deepseekService,
			openaiService,
		},
		aiModel: ai,
	}
}

func (a *AiDispatcherImpl) RequestProcessor(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error) {
	for _, service := range *a.services {
		if service.GetModelType() == m.Type {
			return service.GetAIJsonResponse(c, m, requestId)
		}
	}
	return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}

func (a *AiDispatcherImpl) ProcessCreateModel(ctx context.Context, model models.AIModelInterface) error {
	for _, service := range *a.services {
		if service.GetModelType() == model.GetType() {
			return service.CreateAIModel(ctx, model)
		}
	}
	return errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}

func (a *AiDispatcherImpl) ProcessUpdateModel(ctx context.Context, model models.AIModelInterface) (any, error) {
	for _, service := range *a.services {
		if service.GetModelType() == model.GetType() {
			return service.UpdateAIModel(ctx, model)
		}
	}
	return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}

func (a *AiDispatcherImpl) ProcessModelById(ctx context.Context, id uint) (any, error) {
	model, err := a.aiModel.GetModelById(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	for _, service := range *a.services {
		if service.GetModelType() == model.GetType() {
			return service.GetAIModelById(ctx, id)
		}
	}
	return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}

func (a *AiDispatcherImpl) ProcessAllAIModels(ctx context.Context) []any {
	var models []any
	for _, service := range *a.services {
		aiModels, err := service.GetAllAIModels(ctx)
		if err != nil {
			continue
		}
		models = append(models, aiModels...)
	}
	return models
}
