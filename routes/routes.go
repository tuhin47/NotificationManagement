package routes

import (
	"NotificationManagement/controllers"
	"NotificationManagement/domain"
	"NotificationManagement/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, controller domain.CurlController, keycloakMiddleware *echo.MiddlewareFunc) {
	cg := e.Group("/api/curl", *keycloakMiddleware)

	cg.POST("", controller.CurlHandler, middleware.RequireRoles(RoleCurlCreate))
	cg.GET("/:id", controller.GetCurlRequestByID, middleware.RequireRoles(RoleCurlRead))
	cg.PUT("/:id", controller.UpdateCurlRequest, middleware.RequireRoles(RoleCurlUpdate))
	cg.DELETE("/:id", controller.DeleteCurlRequest, middleware.RequireRoles(RoleCurlDelete))
}

func RegisterLLMRoutes(e *echo.Echo, controller domain.LLMController, keycloakMiddleware *echo.MiddlewareFunc) {
	lg := e.Group("/api/llm", *keycloakMiddleware)

	lg.POST("", controller.CreateLLM, middleware.RequireRoles(RoleLLMCreate))
	lg.GET("/:id", controller.GetLLMByID, middleware.RequireRoles(RoleLLMRead))
	lg.GET("", controller.GetAllLLMs, middleware.RequireRoles(RoleLLMRead))
	lg.PUT("/:id", controller.UpdateLLM, middleware.RequireRoles(RoleLLMUpdate))
	lg.DELETE("/:id", controller.DeleteLLM, middleware.RequireRoles(RoleLLMDelete))
}

func RegisterReminderRoutes(e *echo.Echo, controller domain.ReminderController, keycloakMiddleware *echo.MiddlewareFunc) {
	rg := e.Group("/api/reminder", *keycloakMiddleware)

	rg.POST("", controller.CreateReminder, middleware.RequireRoles(RoleReminderCreate))
	rg.GET("/:id", controller.GetReminderByID, middleware.RequireRoles(RoleReminderRead))
	rg.GET("", controller.GetAllReminders, middleware.RequireRoles(RoleReminderRead))
	rg.PUT("/:id", controller.UpdateReminder, middleware.RequireRoles(RoleReminderUpdate))
	rg.DELETE("/:id", controller.DeleteReminder, middleware.RequireRoles(RoleReminderDelete))
}

func RegisterAIRoutes(e *echo.Echo, controller domain.AIRequestController, keycloakMiddleware *echo.MiddlewareFunc) {
	ai := e.Group("/api/ai", *keycloakMiddleware)

	ai.POST("", controller.CreateAIModel, middleware.RequireRoles(RoleAICreate))
	ai.GET("/:id", controller.GetAIModelByID, middleware.RequireRoles(RoleAIRead))
	ai.GET("", controller.GetAllAIModels, middleware.RequireRoles(RoleAIRead))
	ai.PUT("/:id", controller.UpdateAIModel, middleware.RequireRoles(RoleAIUpdate))
	ai.DELETE("/:id", controller.DeleteAIModel, middleware.RequireRoles(RoleAIDelete))

	ai.POST("/make-request", controller.MakeAIRequestHandler, middleware.RequireRoles(RoleMakeRequest))
}

func RegisterUserRoutes(e *echo.Echo, controller domain.UserController, keycloakMiddleware *echo.MiddlewareFunc) {

}

func RegisterNotificationRoutes(e *echo.Echo, notificationController *controllers.NotificationController, keycloakMiddleware *echo.MiddlewareFunc) {
	e.POST("/api/notify", notificationController.Notify, *keycloakMiddleware)
}

func RegisterTelegramRoutes(e *echo.Echo, controller domain.TelegramController, keycloakMiddleware *echo.MiddlewareFunc) {
	tele := e.Group("/api/telegram", *keycloakMiddleware)
	tele.POST("/verify", controller.VerifyOtp)
}
