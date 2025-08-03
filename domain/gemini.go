package domain

import (
	"NotificationManagement/models"
	"google.golang.org/genai"
)

type GeminiService interface {
	AIService[models.GeminiModel, genai.GenerateContentResponse]
}

type GeminiModelRepository interface {
	Repository[models.GeminiModel, uint]
}
