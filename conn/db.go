package conn

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"time"
)

func NewDB() *gorm.DB {
	gormLogger := NewGormZapLogger(
		logger2.Info,
		200*time.Millisecond,
	)
	dB, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
		panic(err.Error())
	}

	if err := dB.AutoMigrate(
		&models.CurlRequest{},
		&models.Reminder{},
		&models.RequestAIModel{},
		&models.DeepseekModel{},
		&models.AdditionalFields{},
		&models.User{},
		&models.Telegram{},
	); err != nil {
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
		panic(err.Error())
	}
	if err := dB.AutoMigrate(
		&models.GeminiModel{},
	); err != nil {
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
		panic(err.Error())
	}
	log.Info("Database connection successful...")
	return dB
}
