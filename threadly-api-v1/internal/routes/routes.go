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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(httpConfig config.HTTPConfig, pool *pgxpool.Pool, appLogger *zap.Logger) (*gin.Engine, error) {
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

	r.GET("/health", handler.NewHealthHandler(pool, appLogger).Check)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/user/:id", handler.NewUserHandler().GetUserByID)
	}

	return r, nil
}
