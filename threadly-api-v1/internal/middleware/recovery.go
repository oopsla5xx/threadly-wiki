package middleware

import (
	"net/http"
	"threadly-api-v1/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    response.ErrCodeInternalServerError,
					Message: err,
				},
			})
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
}
	