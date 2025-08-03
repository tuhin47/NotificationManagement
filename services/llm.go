package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type LLMServiceImpl struct {
	domain.CommonService[models.RequestAIModel]
	Repo domain.LLMRepository
}

func NewLLMService(repo domain.LLMRepository) domain.LLMService {
	service := &LLMServiceImpl{
		Repo: repo,
	}
	service.CommonService = NewCommonService[models.RequestAIModel](repo, service)
	return service
}
