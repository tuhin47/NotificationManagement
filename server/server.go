package server

import (
	"NotificationManagement/controllers"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/middleware"
	"NotificationManagement/repositories"
	"NotificationManagement/routes"
	"NotificationManagement/services"
	"NotificationManagement/utils/errutil"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(interceptLogger)
	e.Use(errutil.ErrorHandler())
	return e
}

func RegisterRoutes(e *echo.Echo, curlController domain.CurlController, llmController domain.LLMController, reminderController domain.ReminderController, deepseekController domain.AIModelController, aiController domain.AIRequestController) {
	keycloakMiddleware := middleware.KeycloakMiddleware()
	routes.RegisterCurlRoutes(e, curlController, &keycloakMiddleware)
	routes.RegisterLLMRoutes(e, llmController, &keycloakMiddleware)
	routes.RegisterReminderRoutes(e, reminderController, &keycloakMiddleware)
	routes.RegisterDeepseekModelRoutes(e, deepseekController, &keycloakMiddleware)
	routes.RegisterAIRoutes(e, aiController, &keycloakMiddleware)
}

func interceptLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		logger.Info("Incoming request",
			"method", req.Method,
			"path", req.URL.Path,
			"query", req.URL.RawQuery,
			"remote_ip", c.RealIP(),
			"user_agent", req.UserAgent(),
		)
		return next(c)
	}
}

var Module = fx.Options(
	fx.Provide(
		NewEcho,
		controllers.NewAIModelController,
		controllers.NewAIRequestController,
		controllers.NewCurlController,
		controllers.NewLLMController,
		controllers.NewReminderController,

		repositories.NewAIModelRepository,
		repositories.NewCurlRequestRepository,
		repositories.NewAdditionalFieldsRepository,
		repositories.NewDeepseekModelRepository,
		repositories.NewGeminiRepository,
		repositories.NewLLMRepository,
		repositories.NewReminderRepository,

		services.NewAIServiceManager,
		services.NewAIModelService,
		services.NewCurlService,
		services.NewDeepseekModelService,
		services.NewGeminiService,
		services.NewLLMService,
		services.NewReminderService,
	),
	fx.Invoke(RegisterRoutes),
)
