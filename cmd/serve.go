package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"NotificationManagement/server"
	"fmt"
	"net/http"
	"os"

	"context"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the notification management server",
	Long:  `Start the notification management server with the specified configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(
				func() (*gorm.DB, error) {
					db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
					if err != nil {
						logger.Error("Failed to connect to database", "error", err)
						return nil, err
					}
					return db, nil
				},
			),
			server.Module,
			fx.Invoke(func(lc fx.Lifecycle, e *echo.Echo) {
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

				lc.Append(fx.Hook{
					OnStart: func(startCtx context.Context) error {
						go func() {
							if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
								logger.Error("Failed to start server", "error", err)
								os.Exit(1)
							}
						}()
						return nil
					},
					OnStop: func(stopCtx context.Context) error {
						return e.Shutdown(stopCtx)
					},
				})
			},
			),
		)
		app.Run()
	},
}
