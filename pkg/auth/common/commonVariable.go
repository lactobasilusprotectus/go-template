package common

import (
	"fmt"
	"time"
)

var (
	ErrAuthUnauthenticated = fmt.Errorf("unauthenticated")
	ErrEmailNotFound       = fmt.Errorf("email not found")
	ErrWrongPassword       = fmt.Errorf("wrong password")
	ErrUserNotFound        = fmt.Errorf("user not found")
)

const (
	AccessTokenType      = "access_token"
	RefreshTokenType     = "refresh_token"
	AccessTokenLifetime  = time.Minute * 5    // 5 mins
	RefreshTokenLifetime = time.Hour * 24 * 2 // 48 hours

	SessionInvalidated = "1"
)

type LoginToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginInfo struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_uuid"`
}

type LogoutInfo struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,gt=8"`
}
