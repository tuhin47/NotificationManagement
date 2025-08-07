package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/conn"
	"NotificationManagement/logger"
	"NotificationManagement/server"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the notification management server",
	Long:  `Start the notification management server with the specified configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(conn.NewDB),
			conn.ProvideNotifiers(),
			server.Module,
			fx.Invoke(func(lc fx.Lifecycle, e *echo.Echo) {
				// Health check route
				e.GET("/health", func(c echo.Context) error {
					return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
				})

				port := *config.App().Port
				addr := fmt.Sprintf(":%d", port)

				logger.Info("Starting server", "name", config.App().Name, "port", port)
				logger.Info("Server mode", "env", config.App().Env)

				lc.Append(fx.Hook{
					OnStart: func(startCtx context.Context) error {
						go func() {
							if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
