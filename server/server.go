package server

import (
	"NotificationManagement/controllers"
	"NotificationManagement/routes"
	"NotificationManagement/services"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEcho() *echo.Echo {
	return echo.New()
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
