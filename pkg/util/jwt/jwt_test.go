package jwt

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/lactobasilusprotectus/go-template/pkg/common/constant"
	commonTime "github.com/lactobasilusprotectus/go-template/pkg/common/time"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	secret = "super-secret"
)

var (
	now           = time.Now()
	fiveMinsLater = now.Add(5 * time.Minute)
	fiveMinsAgo   = now.Add(-5 * time.Minute)
)

func initMocks(t *testing.T) (mocks map[int]interface{}, deferFunc func()) {
	// mock used modules
	mockCtrl := gomock.NewController(t)
	mockTime := commonTime.NewMockTimeInterface(mockCtrl)

	mocks = map[int]interface{}{
		constant.MockController: mockCtrl,
		constant.MockTime:       mockTime,
	}
	deferFunc = func() {
		mockCtrl.Finish()
	}
	return
}

func initModule(mocks map[int]interface{}, mockFunc func()) (module *JwtModule) {
	mockTime := mocks[constant.MockTime].(*commonTime.MockTimeInterface)
	mockFunc()

	module = New(mockTime)
	module.secret = secret

	return
}

func TestGenerateToken(t *testing.T) {
	testCases := []struct {
		name   string
		doMock func(map[int]interface{})
		err    error
	}{
		{
			name: "positive",
			doMock: func(mocks map[int]interface{}) {
				mockTime := mocks[constant.MockTime].(*commonTime.MockTimeInterface)

				mockTime.EXPECT().Now().Return(time.Now())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// prepare mocked module
			mocks, deferFunc := initMocks(t)
			module := initModule(mocks, func() { tc.doMock(mocks) })
			defer deferFunc()

			// call the function
			token, err := module.GenerateToken(context.Background(), JwtData{
				SessionID:  "session",
				IdentityID: 10,
				Type:       "1",
				Lifetime:   time.Minute,
			})

			// assert returned values
			assert.Equal(t, tc.err, err)
			assert.NotZero(t, token)
		})
	}
}

func TestExtractToken(t *testing.T) {
	testCases := []struct {
		name          string
		generateToken func() string
		data          JwtData
		err           error
	}{
		{
			name: "ParseWithClaims_returns_err",
			generateToken: func() string {
				return "token"
			},
			err: ErrTokenInvalid,
		},
		{
			name: "token_expires",
			generateToken: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(time.Unix(fiveMinsAgo.Unix(), 0)),
						ExpiresAt: jwt.NewNumericDate(time.Unix(fiveMinsAgo.Unix(), 0)),
					},
					SessionID:  "session",
					IdentityID: 10,
					Type:       "type",
				})
				tokenString, _ := token.SignedString([]byte(secret))

				return tokenString
			},
			err: ErrTokenExpired,
		},
		{
			name: "token_not_issued_yet",
			generateToken: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(time.Unix(fiveMinsLater.Unix(), 0)),
						ExpiresAt: jwt.NewNumericDate(time.Unix(fiveMinsLater.Unix(), 0)),
					},
					SessionID:  "session",
					IdentityID: 10,
					Type:       "type",
				})
				tokenString, _ := token.SignedString([]byte(secret))

				return tokenString
			},
			err: ErrTokenInvalid,
		},
		{
			name: "positive",
			generateToken: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(time.Unix(fiveMinsAgo.Unix(), 0)),
						ExpiresAt: jwt.NewNumericDate(time.Unix(fiveMinsLater.Unix(), 0)),
					},
					SessionID:  "session",
					IdentityID: 10,
					Type:       "type",
				})
				tokenString, _ := token.SignedString([]byte(secret))

				return tokenString
			},
			data: JwtData{
				SessionID:  "session",
				IdentityID: 10,
				Type:       "type",
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// prepare mocked module
			mocks, deferFunc := initMocks(t)
			module := initModule(mocks, func() {})
			defer deferFunc()

			// call the function
			data, err := module.ExtractToken(context.Background(), tc.generateToken())

			// assert returned values
			assert.True(t, errors.Is(err, tc.err))
			assert.Equal(t, tc.data, data)
		})
	}
}
