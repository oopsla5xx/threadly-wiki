package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"threadly-api-v1/internal/config"
)

// New builds a console logger in development and a JSON logger otherwise.
// An empty level keeps zap's default for that environment.
func New(environment string, level string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()

	if environment == config.EnvironmentDevelopment {
		zapConfig = zap.NewDevelopmentConfig()
	}

	if level != "" {
		atomicLevel, err := zap.ParseAtomicLevel(level)
		if err != nil {
			return nil, fmt.Errorf("parse log level: %w", err)
		}

		zapConfig.Level = atomicLevel
	}

	return zapConfig.Build()
}

// RequestLogger logs one line per HTTP request, replacing gin.Logger().
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		logger.Info("http request",
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.Int("status", context.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", context.ClientIP()),
		)
	}
}
