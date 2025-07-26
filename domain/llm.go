package domain

import (
	"NotificationManagement/models"

	"github.com/labstack/echo/v4"
)

type LLMService interface {
	CreateLLM(llm *models.UserLLM) error
	GetLLMByID(id uint) (*models.UserLLM, error)
	GetAllLLMs(limit, offset int) ([]models.UserLLM, error)
	UpdateLLM(id uint, llm *models.UserLLM) error
	DeleteLLM(id uint) error
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
