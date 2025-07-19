package cmd

import (
	"NotificationManagement/config"
	"fmt"
	"net/http"
	"os"

	"NotificationManagement/logger"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the notification management server",
	Long:  `Start the notification management server with the specified configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		// Health check route
		e.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		port := config.App().Port
		if port == 0 {
			port = 8080
		}
		addr := fmt.Sprintf(":%d", port)

		logger.Info("Starting server", "name", config.App().Name, "port", port)
		logger.Info("Server mode", "env", config.App().Env)

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	},
}
