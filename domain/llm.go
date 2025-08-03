package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type LLMService interface {
	CommonService[models.RequestAIModel]
}

type LLMRepository interface {
	Repository[models.RequestAIModel, uint]
}

type LLMController interface {
	CreateLLM(c echo.Context) error
	GetLLMByID(c echo.Context) error
	GetAllLLMs(c echo.Context) error
	UpdateLLM(c echo.Context) error
	DeleteLLM(c echo.Context) error
}
