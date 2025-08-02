package routes

import (
	"NotificationManagement/domain"
	"NotificationManagement/middleware"
	"NotificationManagement/utils"

	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, curlController domain.CurlController, keycloakMiddleware *echo.MiddlewareFunc) {
	curlGroup := e.Group("/api/curl", *keycloakMiddleware)

	curlGroup.POST("", curlController.CurlHandler, middleware.RequireRoles(utils.RoleCurlCreate))
	curlGroup.GET("/:id", curlController.GetCurlRequestByID, middleware.RequireRoles(utils.RoleCurlRead))
	curlGroup.PUT("/:id", curlController.UpdateCurlRequest, middleware.RequireRoles(utils.RoleCurlUpdate))
	curlGroup.DELETE("/:id", curlController.DeleteCurlRequest, middleware.RequireRoles(utils.RoleCurlDelete))
}

func RegisterLLMRoutes(e *echo.Echo, llmController domain.LLMController, keycloakMiddleware *echo.MiddlewareFunc) {
	llmGroup := e.Group("/api/llm", *keycloakMiddleware)

	llmGroup.POST("", llmController.CreateLLM, middleware.RequireRoles(utils.RoleLLMCreate))
	llmGroup.GET("/:id", llmController.GetLLMByID, middleware.RequireRoles(utils.RoleLLMRead))
	llmGroup.GET("", llmController.GetAllLLMs, middleware.RequireRoles(utils.RoleLLMRead))
	llmGroup.PUT("/:id", llmController.UpdateLLM, middleware.RequireRoles(utils.RoleLLMUpdate))
	llmGroup.DELETE("/:id", llmController.DeleteLLM, middleware.RequireRoles(utils.RoleLLMDelete))
}

func RegisterReminderRoutes(e *echo.Echo, reminderController domain.ReminderController, keycloakMiddleware *echo.MiddlewareFunc) {
	reminderGroup := e.Group("/api/reminder", *keycloakMiddleware)

	reminderGroup.POST("", reminderController.CreateReminder, middleware.RequireRoles(utils.RoleReminderCreate))
	reminderGroup.GET("/:id", reminderController.GetReminderByID, middleware.RequireRoles(utils.RoleReminderRead))
	reminderGroup.GET("", reminderController.GetAllReminders, middleware.RequireRoles(utils.RoleReminderRead))
	reminderGroup.PUT("/:id", reminderController.UpdateReminder, middleware.RequireRoles(utils.RoleReminderUpdate))
	reminderGroup.DELETE("/:id", reminderController.DeleteReminder, middleware.RequireRoles(utils.RoleReminderDelete))
}

func RegisterDeepseekModelRoutes(e *echo.Echo, controller domain.AIModelController, keycloakMiddleware *echo.MiddlewareFunc) {
	deepseekGroup := e.Group("/api/deepseek-model", *keycloakMiddleware)

	deepseekGroup.POST("", controller.CreateAIModel, middleware.RequireRoles(utils.RoleDeepseekCreate))
	deepseekGroup.GET("/:id", controller.GetAIModelByID, middleware.RequireRoles(utils.RoleDeepseekRead))
	deepseekGroup.GET("", controller.GetAllAIModels, middleware.RequireRoles(utils.RoleDeepseekRead))
	deepseekGroup.PUT("/:id", controller.UpdateAIModel, middleware.RequireRoles(utils.RoleDeepseekUpdate))
	deepseekGroup.DELETE("/:id", controller.DeleteAIModel, middleware.RequireRoles(utils.RoleDeepseekDelete))
}

func RegisterAIRoutes(e *echo.Echo, aiController domain.AIRequestController, keycloakMiddleware *echo.MiddlewareFunc) {
	aiGroup := e.Group("/api/ai", *keycloakMiddleware)

	aiGroup.POST("/make-request", aiController.MakeAIRequestHandler, middleware.RequireRoles(utils.RoleMakeRequest))
}
