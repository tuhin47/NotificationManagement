package domain

import (
	"context"

	"github.com/labstack/echo/v4"
)

type CommonService[T any] interface {
	CreateModel(c echo.Context, entity *T) error
	GetModelById(c echo.Context, id uint) (*T, error)
	GetAllModels(c echo.Context, limit, offset int) ([]T, error)
	UpdateModel(c echo.Context, id uint, model *T) (*T, error)
	DeleteModel(c echo.Context, id uint) error
	GetInstance() CommonService[T]
	GetContext() context.Context
}
