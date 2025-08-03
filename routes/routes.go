package routes

import (
	"NotificationManagement/domain"
	"NotificationManagement/middleware"
	"NotificationManagement/utils"

	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, controller domain.CurlController, keycloakMiddleware *echo.MiddlewareFunc) {
	cg := e.Group("/api/curl", *keycloakMiddleware)

	cg.POST("", controller.CurlHandler, middleware.RequireRoles(utils.RoleCurlCreate))
	cg.GET("/:id", controller.GetCurlRequestByID, middleware.RequireRoles(utils.RoleCurlRead))
	cg.PUT("/:id", controller.UpdateCurlRequest, middleware.RequireRoles(utils.RoleCurlUpdate))
	cg.DELETE("/:id", controller.DeleteCurlRequest, middleware.RequireRoles(utils.RoleCurlDelete))
}

func RegisterLLMRoutes(e *echo.Echo, controller domain.LLMController, keycloakMiddleware *echo.MiddlewareFunc) {
	lg := e.Group("/api/llm", *keycloakMiddleware)

	lg.POST("", controller.CreateLLM, middleware.RequireRoles(utils.RoleLLMCreate))
	lg.GET("/:id", controller.GetLLMByID, middleware.RequireRoles(utils.RoleLLMRead))
	lg.GET("", controller.GetAllLLMs, middleware.RequireRoles(utils.RoleLLMRead))
	lg.PUT("/:id", controller.UpdateLLM, middleware.RequireRoles(utils.RoleLLMUpdate))
	lg.DELETE("/:id", controller.DeleteLLM, middleware.RequireRoles(utils.RoleLLMDelete))
}

func RegisterReminderRoutes(e *echo.Echo, controller domain.ReminderController, keycloakMiddleware *echo.MiddlewareFunc) {
	rg := e.Group("/api/reminder", *keycloakMiddleware)

	rg.POST("", controller.CreateReminder, middleware.RequireRoles(utils.RoleReminderCreate))
	rg.GET("/:id", controller.GetReminderByID, middleware.RequireRoles(utils.RoleReminderRead))
	rg.GET("", controller.GetAllReminders, middleware.RequireRoles(utils.RoleReminderRead))
	rg.PUT("/:id", controller.UpdateReminder, middleware.RequireRoles(utils.RoleReminderUpdate))
	rg.DELETE("/:id", controller.DeleteReminder, middleware.RequireRoles(utils.RoleReminderDelete))
}

func RegisterAIRoutes(e *echo.Echo, controller domain.AIRequestController, keycloakMiddleware *echo.MiddlewareFunc) {
	ai := e.Group("/api/ai", *keycloakMiddleware)

	ai.POST("", controller.CreateAIModel, middleware.RequireRoles(utils.RoleDeepseekCreate))
	ai.GET("/:id", controller.GetAIModelByID, middleware.RequireRoles(utils.RoleDeepseekRead))
	ai.GET("", controller.GetAllAIModels, middleware.RequireRoles(utils.RoleDeepseekRead))
	ai.PUT("/:id", controller.UpdateAIModel, middleware.RequireRoles(utils.RoleDeepseekUpdate))
	ai.DELETE("/:id", controller.DeleteAIModel, middleware.RequireRoles(utils.RoleDeepseekDelete))

	ai.POST("/make-request", controller.MakeAIRequestHandler, middleware.RequireRoles(utils.RoleMakeRequest))
}
