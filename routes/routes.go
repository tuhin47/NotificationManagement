package routes

import (
	"NotificationManagement/domain"
	"NotificationManagement/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, curlController domain.CurlController) {
	curlGroup := e.Group("/api/curl", middleware.KeycloakMiddleware())

	curlGroup.POST("", curlController.CurlHandler, middleware.RequireRoles("curl_create"))
	curlGroup.GET("/:id", curlController.GetCurlRequestByID, middleware.RequireRoles("curl_read"))
	curlGroup.PUT("/:id", curlController.UpdateCurlRequest, middleware.RequireRoles("curl_update"))
	curlGroup.DELETE("/:id", curlController.DeleteCurlRequest, middleware.RequireRoles("curl_delete"))
}

func RegisterLLMRoutes(e *echo.Echo, llmController domain.LLMController) {
	llmGroup := e.Group("/api/llm", middleware.KeycloakMiddleware())

	llmGroup.POST("", llmController.CreateLLM, middleware.RequireRoles("llm_create"))
	llmGroup.GET("/:id", llmController.GetLLMByID, middleware.RequireRoles("llm_read"))
	llmGroup.GET("", llmController.GetAllLLMs, middleware.RequireRoles("llm_read"))
	llmGroup.PUT("/:id", llmController.UpdateLLM, middleware.RequireRoles("llm_update"))
	llmGroup.DELETE("/:id", llmController.DeleteLLM, middleware.RequireRoles("llm_delete"))
}

func RegisterReminderRoutes(e *echo.Echo, reminderController domain.ReminderController) {
	reminderGroup := e.Group("/api/reminder", middleware.KeycloakMiddleware())

	reminderGroup.POST("", reminderController.CreateReminder, middleware.RequireRoles("reminder_create"))
	reminderGroup.GET("/:id", reminderController.GetReminderByID, middleware.RequireRoles("reminder_read"))
	reminderGroup.GET("", reminderController.GetAllReminders, middleware.RequireRoles("reminder_read"))
	reminderGroup.PUT("/:id", reminderController.UpdateReminder, middleware.RequireRoles("reminder_update"))
	reminderGroup.DELETE("/:id", reminderController.DeleteReminder, middleware.RequireRoles("reminder_delete"))
}

func RegisterDeepseekModelRoutes(e *echo.Echo, deepseekController domain.DeepseekModelController) {
	deepseekGroup := e.Group("/api/deepseek-model", middleware.KeycloakMiddleware())

	deepseekGroup.POST("", deepseekController.CreateDeepseekModel, middleware.RequireRoles("deepseek_create"))
	deepseekGroup.GET("/:id", deepseekController.GetDeepseekModelByID, middleware.RequireRoles("deepseek_read"))
	deepseekGroup.GET("", deepseekController.GetAllDeepseekModels, middleware.RequireRoles("deepseek_read"))
	deepseekGroup.PUT("/:id", deepseekController.UpdateDeepseekModel, middleware.RequireRoles("deepseek_update"))
	deepseekGroup.DELETE("/:id", deepseekController.DeleteDeepseekModel, middleware.RequireRoles("deepseek_delete"))
}

func RegisterAIRoutes(e *echo.Echo, aiController domain.AIController) {
	aiGroup := e.Group("/api/ai", middleware.KeycloakMiddleware())

	aiGroup.POST("/make-request", aiController.MakeAIRequestHandler, middleware.RequireRoles("ai_create"))
}
