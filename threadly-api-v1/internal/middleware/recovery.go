package middleware

import (
	"threadly-api-v1/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery turns any panic into the standard 500 response.
// gin.CustomRecovery already logs the panic value with its stack trace.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ any) {
		appErr := response.InternalServerError()
		c.AbortWithStatusJSON(appErr.StatusCode, response.Response{
			Success: false,
			Error: &response.ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	})
}
