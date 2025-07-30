package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type DeepseekModelService interface {
	CreateDeepseekModel(model *models.DeepseekModel) error
	GetDeepseekModelByID(id uint) (*models.DeepseekModel, error)
	GetAllDeepseekModels(limit, offset int) ([]models.DeepseekModel, error)
	UpdateDeepseekModel(id uint, model *models.DeepseekModel) error
	DeleteDeepseekModel(id uint) error
}

type DeepseekModelRepository interface {
	Repository[models.DeepseekModel, uint]
}

type DeepseekModelController interface {
	CreateDeepseekModel(c echo.Context) error
	GetDeepseekModelByID(c echo.Context) error
	GetAllDeepseekModels(c echo.Context) error
	UpdateDeepseekModel(c echo.Context) error
	DeleteDeepseekModel(c echo.Context) error
}
