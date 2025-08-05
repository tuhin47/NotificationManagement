package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type DeepseekService interface {
	AIService[models.DeepseekModel]
	PullModel(c echo.Context, model *models.DeepseekModel) error
}

type DeepseekModelRepository interface {
	Repository[models.DeepseekModel, uint]
}
