package domain

import (
	"NotificationManagement/models"
)

type OpenAIService interface {
	AIService[models.OpenAIModel]
}

type OpenAIModelRepository interface {
	Repository[models.OpenAIModel, uint]
}
