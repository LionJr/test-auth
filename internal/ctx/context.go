package ctx

import (
	"github.com/gin-gonic/gin"
)

const (
	KeyUserID = "user_id"
)

func UserID(ctx *gin.Context) int {
	p, ok := ctx.Get(KeyUserID)
	if !ok {
		return 0
	}

	userID, ok := p.(int)
	if !ok {
		return 0
	}

	return userID
}
