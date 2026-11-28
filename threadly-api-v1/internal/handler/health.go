package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"threadly-api-v1/pkg/response"
)

const healthCheckTimeout = 2 * time.Second

type HealthHandler struct {
	pool        *pgxpool.Pool
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewHealthHandler(pool *pgxpool.Pool, redisClient *redis.Client, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		pool:        pool,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (h *HealthHandler) Check(ginContext *gin.Context) {
	pingContext, cancel := context.WithTimeout(ginContext.Request.Context(), healthCheckTimeout)
	defer cancel()

	if err := h.pool.Ping(pingContext); err != nil {
		h.logger.Error("health check: database ping failed", zap.Error(err))
		writeServiceUnavailable(ginContext)
		return
	}

	if err := h.redisClient.Ping(pingContext).Err(); err != nil {
		h.logger.Error("health check: redis ping failed", zap.Error(err))
		writeServiceUnavailable(ginContext)
		return
	}

	ginContext.JSON(http.StatusOK, response.Response{
		Success: true,
		Data: gin.H{
			"status":   "ok",
			"database": "ok",
			"redis":    "ok",
		},
	})
}

func writeServiceUnavailable(ginContext *gin.Context) {
	appError := response.ServiceUnavailable()

	ginContext.JSON(appError.StatusCode, response.Response{
		Success: false,
		Error: &response.ErrorInfo{
			Code:    appError.Code,
			Message: appError.Message,
		},
	})
}
