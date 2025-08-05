package db

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"time"
)

func NewDB() (*gorm.DB, error) {
	gormLogger := NewGormZapLogger(
		logger2.Info,
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
		&models.User{},
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
