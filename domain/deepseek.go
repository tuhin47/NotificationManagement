package domain

import (
	"NotificationManagement/models"
	"context"
)

type DeepseekService interface {
	AIService[models.DeepseekModel]
	PullModel(c context.Context, model *models.DeepseekModel) error
}

type DeepseekModelRepository interface {
	Repository[models.DeepseekModel, uint]
}
