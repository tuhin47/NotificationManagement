package repositories

import (
	"NotificationManagement/utils/errutil"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SQLRepository[T any] struct {
	db *gorm.DB
}

type Filter struct {
	Field string
	Op    string // e.g., "=", "LIKE", "IN"
	Value interface{}
}

type ContextKey struct {
	Filter *[]Filter
}

func NewSQLRepository[T any](db *gorm.DB) *SQLRepository[T] {
	return &SQLRepository[T]{db: db}
}
func (r *SQLRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	return r.db
}

func (r *SQLRepository[T]) Create(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return handleDbError(err)
	}
	return nil
}

func handleDbError(err error) error {
	//https://www.postgresql.org/docs/current/errcodes-appendix.html
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}
	var target *pgconn.PgError
	if errors.As(err, &target) {
		switch target.Code {
		case "23505":
			return errutil.NewAppErrorWithMessage(errutil.ErrDuplicateEntry, err, target.Detail)
		case "23503":
			return errutil.NewAppErrorWithMessage(errutil.ErrRecordNotFound, err, target.Detail)
		}
	}
	return errutil.NewAppError(errutil.ErrDatabaseQuery, err)
}

func (r *SQLRepository[T]) GetByID(ctx context.Context, id uint, preloads *[]string) (*T, error) {
	var entity T
	dbContext := r.db.WithContext(ctx)
	if preloads != nil {
		for _, it := range *preloads {
			dbContext = dbContext.Preload(it)
		}
	}
	err := dbContext.First(&entity, id).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &entity, nil
}

func (r *SQLRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	withContext := r.db.WithContext(ctx)
	withContext = ApplyFilter(ctx, withContext)
	err := withContext.Limit(limit).Offset(offset).Find(&entities).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return entities, nil
}

func ApplyFilter(ctx context.Context, query *gorm.DB) *gorm.DB {
	key := ContextKey{}

	type ExtraFilters *[]Filter
	if contextKey, ok := ctx.Value(key).(*ContextKey); ok {
		for _, f := range *contextKey.Filter {
			clause := fmt.Sprintf("%s %s ?", f.Field, f.Op)
			query = query.Where(clause, f.Value)
		}
	}
	return query
}

func (r *SQLRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		return handleDbError(err)
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
		return 0, handleDbError(err)
	}
	return count, nil
}

func (r *SQLRepository[T]) GetByIDs(ctx context.Context, ids []uint, preloads *[]string) ([]T, error) {
	var entities []T
	dbContext := r.db.WithContext(ctx)
	if preloads != nil {
		for _, it := range *preloads {
			dbContext = dbContext.Preload(it)
		}
	}
	err := dbContext.Where("id IN (?)", ids).Find(&entities).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return entities, nil
}
