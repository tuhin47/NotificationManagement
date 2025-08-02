package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type LLMServiceImpl struct {
	domain.CommonService[models.UserLLM]
	Repo domain.LLMRepository
}

func NewLLMService(repo domain.LLMRepository) domain.LLMService {
	service := &LLMServiceImpl{
		Repo: repo,
	}
	service.CommonService = NewCommonService[models.UserLLM](repo, service)
	return service
}
