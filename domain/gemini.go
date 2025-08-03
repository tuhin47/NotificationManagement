package domain

import (
	"NotificationManagement/models"
)

type GeminiService interface {
	AIService[models.GeminiModel]
}

type GeminiModelRepository interface {
	Repository[models.GeminiModel, uint]
}
