package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type AIModelServiceImpl struct {
	domain.CommonService[models.AIModel]
}

func NewAIModelService(repo domain.AIModelRepository) domain.AIModelService {
	service := &AIModelServiceImpl{}
	service.CommonService = NewCommonService(repo, service)
	return service
}
