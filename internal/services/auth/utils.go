package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func hashRefreshToken(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
