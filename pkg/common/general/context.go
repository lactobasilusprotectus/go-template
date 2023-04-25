package general

import (
	"context"
	"github.com/lactobasilusprotectus/go-template/pkg/common/constant"
)

func GetUserIDFromCtx(ctx context.Context) (userID int64, ok bool) {
	userID, ok = ctx.Value(constant.ContextKeyUserID).(int64)
	return
}

func SetUserIDIntoCtx(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, constant.ContextKeyUserID, userID)
}

func SetSessionIDIntoCtx(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, constant.ContextKeySession, sessionID)
}
