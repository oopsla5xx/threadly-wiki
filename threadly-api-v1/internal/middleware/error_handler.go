package middleware

import (
	"errors"
	"net/http"
	"threadly-api-v1/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Next()

    if len(c.Errors) == 0 {
      return
    }

    err := c.Errors.Last().Err
    var appErr *response.AppError
    if errors.As(err, &appErr) {
      c.JSON(appErr.StatusCode, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
    } else {
      c.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    response.ErrCodeInternalServerError,
				},
			})
    }
  }
}

