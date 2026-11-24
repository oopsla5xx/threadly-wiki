package routes

import (
	"threadly-api-v1/internal/handler"
	"threadly-api-v1/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
  r := gin.New()
  r.Use(
    middleware.Recovery(),
    middleware.ErrorHandler(),
    cors.New(cors.Config{
      AllowOrigins:     []string{"https://example.com"},
      AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
      AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
      ExposeHeaders:    []string{"Content-Length"},
      AllowCredentials: true,
      MaxAge:           12 * time.Hour,
    }),
  )

  apiV1 := r.Group("/api/v1")
  {
    apiV1.GET("/user/:id", handler.NewUserHandler().GetUserByID)
  }

  return r
}

