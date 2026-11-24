package routes

import (
	"fmt"
	"threadly-api-v1/internal/config"
	"threadly-api-v1/internal/handler"
	"threadly-api-v1/internal/infrastructure/logger"
	"threadly-api-v1/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers holds the HTTP handlers the router mounts; they are built in the composition root.
type Handlers struct {
	Health *handler.HealthHandler
	User   *handler.UserHandler
}

func NewRouter(httpConfig config.HTTPConfig, appLogger *zap.Logger, handlers Handlers) (*gin.Engine, error) {
	r := gin.New()

	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}

	r.Use(
		logger.RequestLogger(appLogger),
		middleware.Recovery(),
		middleware.ErrorHandler(),
		cors.New(cors.Config{
			AllowOrigins:     httpConfig.CORS.AllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	r.GET("/health", handlers.Health.Check)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/user/:id", handlers.User.GetUserByID)
	}

	return r, nil
}
