package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
)

type CommonServiceImpl[T any] struct {
	Repo     domain.Repository[T, uint]
	Instance domain.CommonService[T]
}

func NewCommonService[T any](repo domain.Repository[T, uint], instance domain.CommonService[T]) *CommonServiceImpl[T] {
	return &CommonServiceImpl[T]{Repo: repo, Instance: instance}
}

func (s *CommonServiceImpl[T]) GetContext() context.Context {
	return context.Background()
}
func (s *CommonServiceImpl[T]) GetInstance() domain.CommonService[T] {
	return s.Instance
}
func (s *CommonServiceImpl[T]) CreateModel(entity *T) error {
	return s.Repo.Create(s.Instance.GetContext(), entity)
}

func (s *CommonServiceImpl[T]) GetModelByID(id uint) (*T, error) {
	model, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *CommonServiceImpl[T]) GetAllModels(limit, offset int) ([]T, error) {
	m, err := s.Repo.GetAll(s.Instance.GetContext(), limit, offset)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CommonServiceImpl[T]) UpdateModel(id uint, model *T) (*T, error) {
	// Check if the model implements ModelInterface
	modelUpdater, ok := any(model).(models.ModelInterface)
	if !ok {
		return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
	}

	existing, err := s.Repo.GetByID(s.Instance.GetContext(), id, nil)
	if err != nil {
		return nil, err
	}
	if existingUpdater, ok := any(existing).(models.ModelInterface); ok {
		existingUpdater.UpdateFromModel(modelUpdater)
	}

	err = s.Repo.Update(s.Instance.GetContext(), existing)
	return existing, err
}

func (s *CommonServiceImpl[T]) DeleteModel(id uint) error {
	return s.Repo.Delete(s.Instance.GetContext(), id)
}
