package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
)

type DeepseekService interface {
	AIService[models.DeepseekModel, types.OllamaResponse]
	PullModel(model *models.DeepseekModel) error
}

type DeepseekModelRepository interface {
	Repository[models.DeepseekModel, uint]
}
