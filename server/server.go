package server

import (
	"NotificationManagement/controllers"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
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

func RegisterRoutes(e *echo.Echo, curlController domain.CurlController, llmController domain.LLMController, reminderController domain.ReminderController) {
	routes.RegisterCurlRoutes(e, curlController)
	routes.RegisterLLMRoutes(e, llmController)
	routes.RegisterReminderRoutes(e, reminderController)
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
		services.NewCurlService,
		controllers.NewCurlController,
		repositories.NewCurlRequestRepository,
		services.NewLLMService,
		controllers.NewLLMController,
		repositories.NewLLMRepository,
		services.NewReminderService,
		controllers.NewReminderController,
		repositories.NewReminderRepository,
	),
	fx.Invoke(RegisterRoutes),
)
