package repositories

import (
	"NotificationManagement/utils/errutil"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SQLRepository[T any] struct {
	db *gorm.DB
}

func NewSQLRepository[T any](db *gorm.DB) *SQLRepository[T] {
	return &SQLRepository[T]{db: db}
}

func (r *SQLRepository[T]) Create(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		var target *pgconn.PgError
		//https://www.postgresql.org/docs/current/errcodes-appendix.html
		if errors.As(err, &target) && target.Code == "23505" {
			return errutil.NewAppErrorWithMessage(errutil.ErrDuplicateEntry, err, target.Detail)
		}
		return errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	return nil
}

func (r *SQLRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errutil.NewAppError(errutil.ErrRecordNotFound, err)
		}
		return nil, errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	return &entity, nil
}

func (r *SQLRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	return entities, nil
}

func (r *SQLRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		return errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	return nil
}

func (r *SQLRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	res := r.db.WithContext(ctx).Delete(&entity, id)
	err := res.Error
	if err != nil {
		return errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	if res.RowsAffected < 1 {
		return errutil.NewAppError(errutil.ErrRecordNotFound, gorm.ErrRecordNotFound)
	}

	return nil
}

func (r *SQLRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}
	return count, nil
}
