package middleware

import (
	"errors"
	"threadly-api-v1/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler renders the last error a handler attached via c.Error.
// Errors other than *response.AppError become the standard 500 response.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		var appErr *response.AppError
		if !errors.As(err, &appErr) {
			appErr = response.InternalServerError()
		}

		c.JSON(appErr.StatusCode, response.Response{
			Success: false,
			Error: &response.ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}
}
