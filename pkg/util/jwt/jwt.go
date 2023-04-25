package jwt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	commonTime "github.com/lactobasilusprotectus/go-template/pkg/common/time"
	"time"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

// jwtClaims is the claims used to generate jwt token
type jwtClaims struct {
	jwt.RegisteredClaims
	SessionID  string
	IdentityID int64
	Type       string
	Lifetime   time.Duration
}

// JwtData is the data used to generate jwt token
type JwtData struct {
	SessionID  string        // unique session identifier
	IdentityID int64         // identity identifier (user ID, etc.)
	Type       string        // token type based on usage (access token, refresh token, etc.)
	Lifetime   time.Duration // expected token lifetime
}

type JwtInterface interface {
	// GenerateToken generate jwt token based on given data
	GenerateToken(ctx context.Context, data JwtData) (token string, err error)

	// ExtractToken the reverse of generate: extract
	ExtractToken(ctx context.Context, token string) (data JwtData, err error)
}

// JwtModule is the jwt module
type JwtModule struct {
	secret string
	time   commonTime.TimeInterface
}

// New creates new JwtModule
func New(t commonTime.TimeInterface) *JwtModule {
	secret := config.Global.JwtSecretAccessToken

	if secret == "" {
		secret = generateRandomString()
	}

	return &JwtModule{
		secret: secret,
		time:   t,
	}
}

// GenerateToken generate jwt token based on given data
func (j *JwtModule) GenerateToken(ctx context.Context, data JwtData) (token string, err error) {
	claims := j.newClaims(data.SessionID, data.IdentityID, data.Type, data.Lifetime)

	tokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenUnsigned.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractToken the reverse of generate: extract
func (j *JwtModule) ExtractToken(ctx context.Context, token string) (data JwtData, err error) {
	// parse tokenString into token
	parsedToken, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	// parsing token returns error
	if err != nil {
		// is it validation error?
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) && validationErr.Errors == jwt.ValidationErrorExpired {
			err = fmt.Errorf("%w: %+v", ErrTokenExpired, err)
			return
		}

		// validation passed but error occurs
		err = fmt.Errorf("%w: %+v", ErrTokenInvalid, err)
		return
	}

	if parsedToken == nil {
		err = fmt.Errorf("%w: %s", ErrTokenInvalid, "jwt.ParseWithClaims returns nil")
		return
	}

	// parse token to jwtClaims
	claims, ok := parsedToken.Claims.(*jwtClaims)
	if !ok {
		err = fmt.Errorf("%w: %s", ErrTokenInvalid, "unable to parse token.Claims to *jwtClaims")
		return
	}

	data.IdentityID = claims.IdentityID
	data.SessionID = claims.SessionID
	data.Type = claims.Type

	return data, nil
}

// generateRandomString returns securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err) //fail to start
	}

	return base64.URLEncoding.EncodeToString(b)
}

// newClaims creates new claims with given sessionID, identityID, type and lifetime.
func (j *JwtModule) newClaims(sessionID string, identityID int64, typ string, lifeTime time.Duration) *jwtClaims {
	now := j.time.Now()

	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(now.Unix(), 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(now.Add(lifeTime).Unix(), 0)),
		},
		SessionID:  sessionID,
		IdentityID: identityID,
		Type:       typ,
	}
}
