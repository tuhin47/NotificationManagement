package domain

import (
	"context"
)

type CommonService[T any] interface {
	CreateModel(c context.Context, entity *T) error
	GetModelById(c context.Context, id uint, preloads *[]string) (*T, error)
	GetAllModels(c context.Context, limit, offset int) ([]T, error)
	UpdateModel(c context.Context, id uint, model *T) (*T, error)
	DeleteModel(c context.Context, id uint) error
	GetInstance() CommonService[T]
	ProcessContext(context.Context) context.Context
}
