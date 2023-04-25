package domain

import (
	"context"
	"github.com/gin-gonic/gin"
)

type GinAuthentication interface {
	// JWT authentication
	MustLogin() gin.HandlerFunc
	GetUserIDFromCtx(ctx context.Context) (userID int64, ok bool)
}
