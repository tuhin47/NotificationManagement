package routes

import (
	"NotificationManagement/domain"

	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, curlController domain.CurlController) {
	e.POST("/api/curl", curlController.CurlHandler)
	e.GET("/api/curl/:id", curlController.GetCurlRequestByID)
	e.PUT("/api/curl/:id", curlController.UpdateCurlRequest)
	e.DELETE("/api/curl/:id", curlController.DeleteCurlRequest)
}

func RegisterLLMRoutes(e *echo.Echo, llmController domain.LLMController) {
	e.POST("/api/llm", llmController.CreateLLM)
	e.GET("/api/llm/:id", llmController.GetLLMByID)
	e.GET("/api/llm", llmController.GetAllLLMs)
	e.PUT("/api/llm/:id", llmController.UpdateLLM)
	e.DELETE("/api/llm/:id", llmController.DeleteLLM)
}

func RegisterReminderRoutes(e *echo.Echo, reminderController domain.ReminderController) {
	e.POST("/api/reminder", reminderController.CreateReminder)
	e.GET("/api/reminder/:id", reminderController.GetReminderByID)
	e.GET("/api/reminder", reminderController.GetAllReminders)
	e.PUT("/api/reminder/:id", reminderController.UpdateReminder)
	e.DELETE("/api/reminder/:id", reminderController.DeleteReminder)
}

func RegisterDeepseekModelRoutes(e *echo.Echo, deepseekController domain.DeepseekModelController) {
	e.POST("/api/deepseek-model", deepseekController.CreateDeepseekModel)
	e.GET("/api/deepseek-model/:id", deepseekController.GetDeepseekModelByID)
	e.GET("/api/deepseek-model", deepseekController.GetAllDeepseekModels)
	e.PUT("/api/deepseek-model/:id", deepseekController.UpdateDeepseekModel)
	e.DELETE("/api/deepseek-model/:id", deepseekController.DeleteDeepseekModel)
}

func RegisterAIRoutes(e *echo.Echo, aiController domain.AIController) {
	e.POST("/api/ai/make-request", aiController.MakeAIRequestHandler)
}
