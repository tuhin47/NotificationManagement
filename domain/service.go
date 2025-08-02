package domain

import "context"

type CommonService[T any] interface {
	CreateModel(entity *T) error
	GetModelByID(id uint) (*T, error)
	GetAllModels(limit, offset int) ([]T, error)
	UpdateModel(id uint, model *T) error
	DeleteModel(id uint) error
	GetContext() context.Context
}
