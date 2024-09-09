package middlewares

import (
	"github.com/gin-gonic/gin"
	internal_ctx "test-auth/internal/ctx"
	"test-auth/internal/models"
	"test-auth/pkg/token_manager"

	"net/http"
)

func TokenAuthMiddleware(tokenManager *token_manager.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		claims, err := tokenManager.ParseAccessToken(token)
		if err != nil {
			abort(ctx, err.Error())
			return
		}

		ctx.Set(internal_ctx.KeyUserID, claims.UserID)
	}
}

func abort(ctx *gin.Context, message string) {
	resp := models.ErrorResponse{
		Message: message,
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
}
