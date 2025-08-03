package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
)

type AIModelServiceImpl struct {
	domain.CommonService[models.AIModel]
	Repo domain.AIModelRepository
}

func NewAIModelService(repo domain.AIModelRepository) domain.AIModelService {
	service := &AIModelServiceImpl{
		Repo: repo,
	}
	service.CommonService = NewCommonService[models.AIModel](repo, service)
	return service
}
