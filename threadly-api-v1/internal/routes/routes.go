package routes

import (
	"threadly-api-v1/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
  router := gin.New()

  apiV1 := router.Group("/api/v1")
  {
    apiV1.GET("/user/:id", handler.NewUserHandler().GetUserByID)
  }

  return router
}
