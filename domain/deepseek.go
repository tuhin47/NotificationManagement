package domain

import (
	"NotificationManagement/models"
)

type DeepseekService interface {
	AIService[models.DeepseekModel]
	PullModel(model *models.DeepseekModel) error
}

type DeepseekModelRepository interface {
	Repository[models.DeepseekModel, uint]
}
