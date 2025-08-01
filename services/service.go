package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
)

type CommonServiceImpl[T any] struct {
	Repo domain.Repository[T, uint]
}

func NewCommonService[T any](repo domain.Repository[T, uint]) *CommonServiceImpl[T] {
	return &CommonServiceImpl[T]{Repo: repo}
}

func (s *CommonServiceImpl[T]) CreateModel(entity *T) error {
	return s.Repo.Create(context.Background(), entity)
}

func (s *CommonServiceImpl[T]) GetModelByID(id uint) (*T, error) {
	model, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *CommonServiceImpl[T]) GetAllModels(limit, offset int) ([]T, error) {
	m, err := s.Repo.GetAll(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CommonServiceImpl[T]) UpdateModel(id uint, model *T) error {
	// Check if the model implements ModelInterface
	modelUpdater, ok := any(model).(models.ModelInterface)
	if !ok {
		return errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
	}

	existing, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return err
	}
	if existingUpdater, ok := any(existing).(models.ModelInterface); ok {
		existingUpdater.UpdateFromModel(modelUpdater)
	}

	return s.Repo.Update(context.Background(), existing)
}

func (s *CommonServiceImpl[T]) DeleteModel(id uint) error {
	return s.Repo.Delete(context.Background(), id)
}
