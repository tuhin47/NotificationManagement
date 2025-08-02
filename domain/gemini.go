package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
)

type GeminiService interface {
	AIService[models.GeminiModel, types.GeminiResponse]
}

type GeminiModelRepository interface {
	Repository[models.GeminiModel, uint]
}
