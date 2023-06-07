package domain

import (
	"context"
	"github.com/lactobasilusprotectus/go-template/pkg/auth/common"
)

type AuthUseCase interface {
	Register(user User) (err error)
	Login(ctx context.Context, email, password string) (token common.LoginToken, err error)
	Info(ctx context.Context) (info common.LoginInfo, err error)
	SendEmail(ctx context.Context, request common.LoginRequest) (err error)
}
