package routes

import (
	"NotificationManagement/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterCurlRoutes(e *echo.Echo, curlController *controllers.CurlController) {
	e.POST("/api/curl", curlController.CurlHandler)
}
