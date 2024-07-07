package logger

import (
	"os"
	"strings"

	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewLogger)
}

var logger *Logger

type Logger struct {
	*zap.Logger
}

func NewLogger(env *env.Env) *Logger {
	level := env.GetEnv("LOG_LEVEL")
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}
	config := zap.NewProductionEncoderConfig()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)
	l := zap.New(core)
	logger = &Logger{
		Logger: l,
	}

	return logger
}
