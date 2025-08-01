package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"
)

type AIModelServiceImpl[T domain.AIModelType, Y any] struct {
	Repo domain.Repository[T, uint]
}

func NewAIModelService[T domain.AIModelType, Y any](repo domain.Repository[T, uint]) *AIModelServiceImpl[T, Y] {
	return &AIModelServiceImpl[T, Y]{Repo: repo}
}

func (s *AIModelServiceImpl[T, Y]) CreateModel(model *T) error {
	return s.Repo.Create(context.Background(), model)
}

func (s *AIModelServiceImpl[T, Y]) GetModelByID(id uint) (*T, error) {
	model, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *AIModelServiceImpl[T, Y]) GetAllAIModels(limit, offset int) ([]T, error) {
	m, err := s.Repo.GetAll(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AIModelServiceImpl[T, Y]) UpdateAIModel(id uint, model *T) error {
	existing, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return err
	}

	// Cast to the updater interface
	if existingUpdater, ok := any(existing).(models.AIModelInterface); ok {
		if modelUpdater, next := any(model).(models.AIModelInterface); next {
			existingUpdater.UpdateFromModel(modelUpdater)
		}
	}

	return s.Repo.Update(context.Background(), existing)
}

func (s *AIModelServiceImpl[T, Y]) DeleteAIModel(id uint) error {
	return s.Repo.Delete(context.Background(), id)
}
