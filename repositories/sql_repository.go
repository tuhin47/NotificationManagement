package repositories

import (
	"context"

	"gorm.io/gorm"
)

type SQLRepository[T any] struct {
	db *gorm.DB
}

func NewSQLRepository[T any](db *gorm.DB) *SQLRepository[T] {
	return &SQLRepository[T]{db: db}
}

func (r *SQLRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *SQLRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SQLRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error
	return entities, err
}

func (r *SQLRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *SQLRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *SQLRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(new(T)).Count(&count).Error
	return count, err
}
