package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type AIServiceImpl[T any] struct {
	domain.CommonService[T]
	Instant domain.AIService[T]
}

func NewAIService[T any](repo domain.Repository[T, uint], instance domain.AIService[T]) domain.AIService[T] {
	return &AIServiceImpl[T]{
		CommonService: NewCommonService(repo, instance),
		Instant:       instance,
	}
}

func (A AIServiceImpl[T]) MakeAIRequest(m *models.AIModel, requestId uint) (interface{}, error) {
	return A.Instant.MakeAIRequest(m, requestId)
}
