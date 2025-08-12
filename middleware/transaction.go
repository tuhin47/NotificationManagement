package middleware

import (
	"NotificationManagement/repositories"
	"NotificationManagement/utils/errutil"
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func TransactionMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tx := db.Begin()
			if tx.Error != nil {
				return tx.Error
			}

			// Store the transaction in the context
			ctx := GetContext(c, tx)
			c.SetRequest(c.Request().WithContext(ctx))

			// Call the next handler
			err := next(c)

			if err != nil {
				tx.Rollback()
				return errutil.HandleError(c, err)
			}

			tx.Commit()
			return nil
		}
	}
}

func GetContext(c echo.Context, tx *gorm.DB) context.Context {
	var ctx context.Context
	if c != nil {
		ctx = c.Request().Context()
	} else {
		ctx = context.Background()
	}
	var filters []*repositories.Filter
	return context.WithValue(ctx, repositories.TXContextKey, &repositories.TxContextKey{DB: tx, Filter: filters})
}
