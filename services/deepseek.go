package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"
)

type DeepseekModelServiceImpl struct {
	Repo domain.DeepseekModelRepository
}

func NewDeepseekModelService(repo domain.DeepseekModelRepository) domain.DeepseekModelService {
	return &DeepseekModelServiceImpl{Repo: repo}
}

func (s *DeepseekModelServiceImpl) CreateDeepseekModel(model *models.DeepseekModel) error {
	return s.Repo.Create(context.Background(), model)
}

func (s *DeepseekModelServiceImpl) GetDeepseekModelByID(id uint) (*models.DeepseekModel, error) {
	model, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *DeepseekModelServiceImpl) GetAllDeepseekModels(limit, offset int) ([]models.DeepseekModel, error) {
	models, err := s.Repo.GetAll(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *DeepseekModelServiceImpl) UpdateDeepseekModel(id uint, model *models.DeepseekModel) error {
	existing, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	existing.Name = model.Name
	existing.ModelName = model.ModelName
	existing.ModifiedAt = model.ModifiedAt
	existing.Size = model.Size

	return s.Repo.Update(context.Background(), existing)
}

func (s *DeepseekModelServiceImpl) DeleteDeepseekModel(id uint) error {
	return s.Repo.Delete(context.Background(), id)
}
