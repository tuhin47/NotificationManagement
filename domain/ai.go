package domain

import (
	"NotificationManagement/types"
)

type AIService[T any, Y any] interface {
	MakeAIRequest(request *T) (Y, error)
}

type OllamaService interface {
	AIService[types.OllamaRequest, types.OllamaResponse]
	PullModel(types.OllamaPullRequest) error
}
