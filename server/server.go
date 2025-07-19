package server

import (
	"NotificationManagement/controllers"
	"NotificationManagement/logger"
	"NotificationManagement/routes"
	"NotificationManagement/services"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
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
	})
	return e
}

func RegisterRoutes(e *echo.Echo, curlController *controllers.CurlController) {
	routes.RegisterCurlRoutes(e, curlController)
}

var Module = fx.Options(
	fx.Provide(
		NewEcho,
		services.NewCurlServiceImpl,
		controllers.NewCurlController,
	),
	fx.Invoke(RegisterRoutes),
)
