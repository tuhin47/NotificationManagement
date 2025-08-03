package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type LLMServiceImpl struct {
	domain.CommonService[models.RequestAIModel]
}

func NewLLMService(repo domain.LLMRepository) domain.LLMService {
	service := &LLMServiceImpl{}
	service.CommonService = NewCommonService(repo, service)
	return service
}
