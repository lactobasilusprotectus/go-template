package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lactobasilusprotectus/go-template/pkg/auth/common"
	"github.com/lactobasilusprotectus/go-template/pkg/common/general"
	"github.com/lactobasilusprotectus/go-template/pkg/domain"
	"github.com/lactobasilusprotectus/go-template/pkg/util/cache"
	httputil "github.com/lactobasilusprotectus/go-template/pkg/util/http"
	"github.com/lactobasilusprotectus/go-template/pkg/util/jwt"
)

func (a *AuthUseCase) MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get token from request
		token := general.GetTokenFromRequest(c)

		if token == "" {
			httputil.WriteUnauthorizedResponse(c)
			return
		}

		// Extract and validate token
		tokenValid, user, sessionID, err := a.extractAndValidateToken(ctx, token, common.AccessTokenType)

		if err != nil {
			httputil.WriteUnauthorizedResponse(c)
			return
		}

		if !tokenValid {
			httputil.WriteUnauthorizedResponse(c)
			return
		}

		// write session information into context
		ctx = general.SetUserIDIntoCtx(ctx, user.ID)      // int64
		ctx = general.SetSessionIDIntoCtx(ctx, sessionID) // string

		newReq := c.Request.WithContext(ctx)
		c.Request = newReq

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		c.Next()
	}
}

func (a *AuthUseCase) extractAndValidateToken(ctx context.Context, token string, tokenType string) (valid bool,
	user domain.User, sessionID string, err error) {
	// initially, it is invalid
	valid = false

	// extract token
	jwtData, err := a.jwtModule.ExtractToken(ctx, token)

	if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalid) {
		// token has expired or it is invalid
		err = nil
		return
	}
	if err != nil {
		// other error
		err = fmt.Errorf("extract token error: %+v", err)
		return
	}

	// validate data parsed from token
	if jwtData.IdentityID == 0 || jwtData.SessionID == "" || jwtData.Type != tokenType {
		return
	}
	userID := jwtData.IdentityID

	// check if session has been invalidated
	invalidated, err := a.isSessionInvalidated(jwtData.SessionID)
	if err != nil {
		err = fmt.Errorf("session checking error: %+v", err)
		return
	}
	if invalidated {
		return
	}

	// get user from repository
	users, err := a.userRepo.FindUserByID(userID)
	if err != nil {
		err = fmt.Errorf("GetUserByIDs error: %+v", err)
		return
	}

	valid = true
	user = users
	sessionID = jwtData.SessionID
	return
}

func (a *AuthUseCase) GetUserIDFromCtx(ctx context.Context) (userID int64, ok bool) {
	return general.GetUserIDFromCtx(ctx)
}

// check for invalidated token in cache
func (a *AuthUseCase) isSessionInvalidated(sessionID string) (invalidated bool, err error) {
	cacheKey := invalidSessionCacheKey(sessionID)
	cacheVal, err := cache.String(a.redis.Get(cacheKey))
	if cacheVal == common.SessionInvalidated {
		return true, nil
	}

	// cache return err other than ErrNilReturned
	if !errors.Is(err, cache.ErrNilReturned) {
		return false, fmt.Errorf("isSessionInvalidated err: %+v", err)
	}

	return false, nil
}

func invalidSessionCacheKey(token string) string {
	return fmt.Sprintf("session-invalid:%s", token)
}
