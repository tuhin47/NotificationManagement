package server

import (
	"NotificationManagement/controllers"
	"NotificationManagement/logger"
	"NotificationManagement/middleware"
	"NotificationManagement/repositories"
	"NotificationManagement/routes"
	"NotificationManagement/services"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(interceptLogger)
	e.Use(middleware.ErrorHandler())
	return e
}

func RegisterRoutes(e *echo.Echo, curlController *controllers.CurlController) {
	routes.RegisterCurlRoutes(e, curlController)
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
		services.NewCurlServiceImpl,
		controllers.NewCurlController,
		repositories.NewCurlRequestRepository,
	),
	fx.Invoke(RegisterRoutes),
)
