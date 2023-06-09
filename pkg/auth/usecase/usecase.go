package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lactobasilusprotectus/go-template/pkg/auth/common"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"github.com/lactobasilusprotectus/go-template/pkg/common/password"
	commonTime "github.com/lactobasilusprotectus/go-template/pkg/common/time"
	"github.com/lactobasilusprotectus/go-template/pkg/domain"
	"github.com/lactobasilusprotectus/go-template/pkg/util/jwt"
	"github.com/lactobasilusprotectus/go-template/pkg/util/queue"
	"github.com/lactobasilusprotectus/go-template/pkg/util/redis"
	"time"
)

type AuthUseCase struct {
	userRepo  domain.UserRepository
	jwtModule jwt.JwtInterface
	redis     redis.Interface
	time      commonTime.TimeInterface
	config    config.Config
	client    queue.Interface
}

func NewAuthUseCase(userRepo domain.UserRepository,
	jwtModule jwt.JwtInterface, redis redis.Interface, time commonTime.TimeInterface,
	config config.Config, client queue.Interface) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		jwtModule: jwtModule,
		redis:     redis,
		time:      time,
		config:    config,
		client:    client,
	}
}

func (a *AuthUseCase) Register(user domain.User) (err error) {
	//hash password
	hashedPassword, err := password.HashPassword(user.Password)

	if err != nil {
		return fmt.Errorf("something wrong: %w", err)
	}

	user.Password = hashedPassword

	//save to database
	return a.userRepo.InsertUser(user)
}

func (a *AuthUseCase) Login(ctx context.Context, email, pass string) (token common.LoginToken, err error) {
	//get user from database
	user, err := a.userRepo.FindUserByEmail(email)

	if err != nil {
		return common.LoginToken{}, common.ErrEmailNotFound
	}

	//check password
	if len(user.Email) > 0 {
		if !password.CheckPasswordHash(pass, user.Password) {
			err = common.ErrWrongPassword
			return
		}

		token, err = a.generateLoginToken(ctx, user.ID)
		if err != nil {
			return
		}
	} else {
		err = common.ErrUserNotFound
		return
	}

	return
}

func (a *AuthUseCase) Info(ctx context.Context) (info common.LoginInfo, err error) {
	//TODO implement me
	panic("implement me")
}

// generate new login token with new session ID
func (a *AuthUseCase) generateLoginToken(ctx context.Context, userID int64) (token common.LoginToken, err error) {
	sessionID := uuid.New().String()

	at, err := a.generateToken(ctx, sessionID, userID, common.AccessTokenType, common.AccessTokenLifetime)
	if err != nil {
		return
	}

	rt, err := a.generateToken(ctx, sessionID, userID, common.RefreshTokenType, common.RefreshTokenLifetime)
	if err != nil {
		return
	}

	token.AccessToken = at
	token.RefreshToken = rt
	return
}

// generate new token
func (a *AuthUseCase) generateToken(ctx context.Context, sessionID string, userID int64, tokenType string,
	lifeTime time.Duration) (token string, err error) {
	return a.jwtModule.GenerateToken(ctx, jwt.JwtData{
		SessionID:  sessionID,
		IdentityID: userID,
		Type:       tokenType,
		Lifetime:   lifeTime,
	})
}

func (a *AuthUseCase) SendEmail(ctx context.Context, request common.LoginRequest) (err error) {
	emailDeliveryTask, err := a.EmailDeliveryTask(ctx, request)

	if err != nil {
		return err
	}

	info, err := a.client.EnqueueTask(emailDeliveryTask)

	if err != nil {
		return fmt.Errorf("something wrong: %w", err)
	}

	fmt.Println(info)

	//close connection
	err = a.client.Close()

	if err != nil {
		return fmt.Errorf("something wrong: %w", err)
	}

	return
}
