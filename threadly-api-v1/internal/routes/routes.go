package routes

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine{
  router := gin.New()

  apiV1 := router.Group("/api/v1")

  apiV1.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  return router
}
