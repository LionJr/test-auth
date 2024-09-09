package auth

import (
	"github.com/gin-gonic/gin"
	"test-auth/internal/models"
)

func sendErrorResponse(ctx *gin.Context, msg string, status int) {
	resp := models.ErrorResponse{
		Message: msg,
	}

	ctx.JSON(status, resp)
}

func sendSuccessResponse(ctx *gin.Context, data interface{}, status int) {
	ctx.JSON(status, data)
}
