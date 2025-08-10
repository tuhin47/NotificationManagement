package domain

import (
	"context"
	"gorm.io/gorm"
)

type Repository[T any, ID comparable] interface {
	GetDB(ctx context.Context) *gorm.DB
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id ID, preloads *[]string) (*T, error)
	GetByIDs(ctx context.Context, ids []uint, preloads *[]string) ([]T, error)
	GetAll(ctx context.Context, limit, offset int) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	Count(ctx context.Context) (int64, error)
}
