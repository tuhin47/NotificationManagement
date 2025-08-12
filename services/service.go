package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
)

type CommonServiceImpl[T any] struct {
	domain.Repository[T, uint]
	Instance domain.CommonService[T]
}

func NewCommonService[T any](repo domain.Repository[T, uint], instance domain.CommonService[T]) domain.CommonService[T] {
	return &CommonServiceImpl[T]{Repository: repo, Instance: instance}
}

func (s *CommonServiceImpl[T]) ProcessContext(c context.Context) context.Context {
	return c
}
func (s *CommonServiceImpl[T]) GetInstance() domain.CommonService[T] {
	return s.Instance
}

func (s *CommonServiceImpl[T]) CreateModel(c context.Context, entity *T) error {
	return s.Create(s.Instance.ProcessContext(c), entity)
}

func (s *CommonServiceImpl[T]) GetModelById(c context.Context, id uint, preloads *[]string) (*T, error) {
	model, err := s.GetByID(s.Instance.ProcessContext(c), id, preloads)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *CommonServiceImpl[T]) GetAllModels(c context.Context, limit, offset int) ([]T, error) {
	m, err := s.GetAll(s.Instance.ProcessContext(c), limit, offset)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CommonServiceImpl[T]) UpdateModel(c context.Context, id uint, model *T) (*T, error) {
	modelUpdater, ok := any(model).(models.ModelInterface)
	if !ok {
		return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
	}

	existing, err := s.GetByID(s.Instance.ProcessContext(c), id, nil)
	if err != nil {
		return nil, err
	}
	if existingUpdater, ok := any(existing).(models.ModelInterface); ok {
		existingUpdater.UpdateFromModel(modelUpdater)
	}

	err = s.Update(s.Instance.ProcessContext(c), existing)
	return existing, err
}

func (s *CommonServiceImpl[T]) DeleteModel(c context.Context, id uint) error {
	return s.Delete(s.Instance.ProcessContext(c), id)
}
