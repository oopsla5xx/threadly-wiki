package handler

import (
	"net/http"
	"threadly-api-v1/internal/service"
	"threadly-api-v1/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUserByID(context *gin.Context) {
	id := context.Param("id")
	
	result, _ := h.userService.GetUserByID(id)
	// if err != nil {
	// 	if errors.Is(err, user.ErrNotFound) {
	// 		_ = context.Error(response.NewAppError(
	// 			http.StatusNotFound,
	// 			"USER_NOT_FOUND",
	// 			"User not found",
	// 		))
	// 		return
	// 	}

	// 	_ = context.Error(err)
	// 	return
	// }
	context.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    result,
	})
}

