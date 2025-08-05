package db

import (
	"NotificationManagement/logger"
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type GormZapLogger struct {
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
	ZLogger       *zap.Logger
}

func NewGormZapLogger(level gormlogger.LogLevel, slowThreshold time.Duration) *GormZapLogger {
	return &GormZapLogger{
		LogLevel:      level,
		SlowThreshold: slowThreshold,
		ZLogger:       logger.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(3)), // disables zap's caller annotation
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZLogger.Sugar().Infof(msg, data...)
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}
	duration := time.Since(begin)
	sql, rows := fc()
	msg := "GORM SQL: " + sql
	fields := []interface{}{
		"duration", duration,
		"rows", rows,
	}
	if err != nil && l.LogLevel >= gormlogger.Error {
		l.ZLogger.Sugar().Errorw(msg, append(fields, "error", err)...)
	} else if l.SlowThreshold != 0 && duration > l.SlowThreshold && l.LogLevel >= gormlogger.Warn {
		l.ZLogger.Sugar().Warnw("GORM SLOW SQL: "+sql, append(fields, "slow_threshold", l.SlowThreshold)...)
	} else if l.LogLevel >= gormlogger.Info {
		l.ZLogger.Sugar().Infow(msg, fields...)
	}
}
