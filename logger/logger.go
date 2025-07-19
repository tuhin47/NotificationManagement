package logger

import (
	"NotificationManagement/config"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

// Init initializes the global Zap logger using config.Logger(). Only runs once.
func Init() {
	loggerOnce.Do(func() {
		cfg := config.Logger()
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}

		var encoder zapcore.Encoder
		format := cfg.Format
		if format == "" {
			format = "console"
		}
		if format == "json" {
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		} else {
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		}

		var writers []zapcore.WriteSyncer
		writers = append(writers, zapcore.AddSync(os.Stdout))
		if cfg.FilePath != "" {
			dir := ""
			if idx := len(cfg.FilePath) - len("/app.log"); idx > 0 {
				dir = cfg.FilePath[:idx]
			}
			if dir != "" {
				_ = os.MkdirAll(dir, 0755)
			}
			file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				writers = append(writers, zapcore.AddSync(file))
			}
		}
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(writers...),
			level,
		)
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

// L returns the global Zap logger. Panics if not initialized.
func L() *zap.Logger {
	if logger == nil {
		panic("logger not initialized: call logger.Init() first")
	}
	return logger
}

// Sync flushes any buffered log entries.
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// toZapField converts various types to zap.Field for flexible logging.
func toZapField(val interface{}) zap.Field {
	switch v := val.(type) {
	case zap.Field:
		return v
	case error:
		return zap.Error(v)
	case string:
		return zap.String("msg", v)
	case map[string]interface{}:
		return zap.Any("data", v)
	case map[string]string:
		return zap.Any("data", v)
	// Note: []zap.Field is not supported directly; flatten in the caller if needed.
	default:
		return zap.Any("data", v)
	}
}

// Info logs a message at InfoLevel. Accepts zap.Field, error, map, key-value pairs, or any value.
func Info(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().Info(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().Info(msg, zapFields...)
}

// Error logs a message at ErrorLevel. Accepts zap.Field, error, map, key-value pairs, or any value.
func Error(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().Error(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().Error(msg, zapFields...)
}

// Warn logs a message at WarnLevel. Accepts zap.Field, error, map, key-value pairs, or any value.
func Warn(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().Warn(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().Warn(msg, zapFields...)
}

// Debug logs a message at DebugLevel. Accepts zap.Field, error, map, key-value pairs, or any value.
func Debug(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().Debug(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().Debug(msg, zapFields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1). Accepts zap.Field, error, map, key-value pairs, or any value.
func Fatal(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().Fatal(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().Fatal(msg, zapFields...)
}

// DPanic logs a message at DPanicLevel. In development, it then panics. Accepts zap.Field, error, map, key-value pairs, or any value.
func DPanic(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	if len(fields) > 1 && len(fields)%2 == 0 {
		allStringKeys := true
		for i := 0; i < len(fields); i += 2 {
			if _, ok := fields[i].(string); !ok {
				allStringKeys = false
				break
			}
		}
		if allStringKeys {
			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				val := fields[i+1]
				zapFields = append(zapFields, zap.Any(key, val))
			}
			L().DPanic(msg, zapFields...)
			return
		}
	}
	for _, f := range fields {
		switch v := f.(type) {
		case []zap.Field:
			zapFields = append(zapFields, v...)
		default:
			zapFields = append(zapFields, toZapField(v))
		}
	}
	L().DPanic(msg, zapFields...)
}

// allStrings returns true if all elements are strings.
func allStrings(fields []interface{}) bool {
	for _, f := range fields {
		if _, ok := f.(string); !ok {
			return false
		}
	}
	return true
}
