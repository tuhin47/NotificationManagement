package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type LLMService interface {
	CommonService[models.UserLLM]
}

type LLMRepository interface {
	Repository[models.UserLLM, uint]
}

type LLMController interface {
	CreateLLM(c echo.Context) error
	GetLLMByID(c echo.Context) error
	GetAllLLMs(c echo.Context) error
	UpdateLLM(c echo.Context) error
	DeleteLLM(c echo.Context) error
}
