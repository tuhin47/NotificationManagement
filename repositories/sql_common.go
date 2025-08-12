package repositories

import (
	"NotificationManagement/domain"
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
type txContextKey string

const TXContextKey txContextKey = "txContext"

type TxContextKey struct {
	DB     *gorm.DB
	Filter []*Filter
}

type Filter struct {
	Field   string
	Op      string
	Value   interface{}
	Applied bool
}

func NewFilter(field string, op string, value interface{}) *Filter {
	return &Filter{Value: value, Op: op, Field: field, Applied: false}
}

func NewSQLRepository[T any](db *gorm.DB) domain.Repository[T, uint] {
	return &SQLRepository[T]{db: db}
}

func (r *SQLRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := GetTxContext(ctx); ok {
		return tx.DB.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func GetTxContext(ctx context.Context) (*TxContextKey, bool) {
	if contextKey, ok := ctx.Value(TXContextKey).(*TxContextKey); ok {
		return contextKey, true
	}
	return nil, false
}

func ApplyFilter(ctx context.Context, query *gorm.DB) *gorm.DB {
	if contextKey, ok := GetTxContext(ctx); ok {
		for i, f := range contextKey.Filter {
			if f.Applied {
				continue
			}
			clause := fmt.Sprintf("%s %s ?", f.Field, f.Op)
			query = query.Where(clause, f.Value)
			(contextKey.Filter)[i].Applied = true
		}
	}
	return query
}

func (r *SQLRepository[T]) WithTx(tx *gorm.DB) domain.Repository[T, uint] {
	return &SQLRepository[T]{db: tx}
}

func (r *SQLRepository[T]) Create(ctx context.Context, entity *T) error {
	withContext := r.GetDB(ctx)
	withContext = ApplyFilter(ctx, withContext)
	err := withContext.Create(entity).Error
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
	dbContext := r.GetDB(ctx)
	if preloads != nil {
		for _, it := range *preloads {
			dbContext = dbContext.Preload(it)
		}
	}
	dbContext = ApplyFilter(ctx, dbContext)
	err := dbContext.First(&entity, id).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &entity, nil
}

func (r *SQLRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	withContext := r.GetDB(ctx)
	withContext = ApplyFilter(ctx, withContext)
	err := withContext.Limit(limit).Offset(offset).Find(&entities).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return entities, nil
}

func (r *SQLRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.GetDB(ctx).Save(entity).Error
	if err != nil {
		return handleDbError(err)
	}
	return nil
}

func (r *SQLRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	res := r.GetDB(ctx).Delete(&entity, id)
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
	err := r.GetDB(ctx).Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, handleDbError(err)
	}
	return count, nil
}

func (r *SQLRepository[T]) GetByIDs(ctx context.Context, ids []uint, preloads *[]string) ([]T, error) {
	var entities []T
	dbContext := r.GetDB(ctx)
	if preloads != nil {
		for _, it := range *preloads {
			dbContext = dbContext.Preload(it)
		}
	}
	dbContext = ApplyFilter(ctx, dbContext)
	err := dbContext.Where("id IN (?)", ids).Find(&entities).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return entities, nil
}
