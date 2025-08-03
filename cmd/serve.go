package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/server"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"net/http"
	"os"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the notification management server",
	Long:  `Start the notification management server with the specified configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(NewDB),
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

func NewDB() (*gorm.DB, error) {
	gormLogger := logger.NewGormZapLogger(
		gormlogger.Info,
		200*time.Millisecond,
	)
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
		return nil, err
	}
	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.CurlRequest{},
		&models.Reminder{},
		&models.RequestAIModel{},
		&models.DeepseekModel{},
		&models.AdditionalFields{},
	); err != nil {
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
		return nil, err
	}
	// As grom doesn't support single table model. We need to execute it separately
	if err := db.AutoMigrate(
		&models.GeminiModel{},
	); err != nil {
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
		return nil, err
	}
	return db, nil
}
