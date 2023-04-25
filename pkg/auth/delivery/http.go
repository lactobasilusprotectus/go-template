package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lactobasilusprotectus/go-template/pkg/auth/common"
	"github.com/lactobasilusprotectus/go-template/pkg/domain"
	_ "github.com/lactobasilusprotectus/go-template/pkg/util/http"
	httputil "github.com/lactobasilusprotectus/go-template/pkg/util/http"
)

type AuthHttpHandler struct {
	authMiddleware domain.GinAuthentication
	authUseCase    domain.AuthUseCase
}

func NewAuthHttpHandler(authMiddleware domain.GinAuthentication, authUseCase domain.AuthUseCase) *AuthHttpHandler {
	return &AuthHttpHandler{
		authMiddleware: authMiddleware,
		authUseCase:    authUseCase,
	}
}

func (a *AuthHttpHandler) Register(g *gin.Engine) {
	g.POST("login", a.Login)
	g.POST("register", a.Regis)
}

// Login				godoc
//
//	@Summary		Login user to get token.
//	@Description	Login user to get token.
//	@Produce		application/json
//	@Tags			auth
//	@Param			body	body		common.LoginRequest	true	"Login Request"
//	@Success		200		{object}	http.BaseResponse
//	@Failure		400		{object}	http.BaseResponse
//	@Failure		500		{object}	http.BaseResponse
//	@Router			/login [post]
func (a *AuthHttpHandler) Login(c *gin.Context) {
	// init request body
	var loginRequest common.LoginRequest

	//bind request body
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		httputil.WriteBadRequestResponseWithErrMsg(c, httputil.ResponseBadRequestError, err)
		return
	}

	// validate request body
	if err := validator.New().Struct(&loginRequest); err != nil {
		httputil.WriteBadRequestResponseWithErrMsg(c, httputil.ResponseBadRequestError, err)
		return
	}

	// call use case
	token, err := a.authUseCase.Login(c, loginRequest.Email, loginRequest.Password)

	// handle error
	if err != nil {
		httputil.WriteServerErrorResponse(c, httputil.ResponseServerError, err)
		return
	}

	// write response
	httputil.WriteOkResponse(c, token)
	return
}

// Regis				godoc
//
//	@Summary		Regis user.
//	@Description	Regis user.
//	@Produce		application/json
//	@Tags			auth
//	@Param			body	body		common.RegisterRequest	true	"Registration Request"
//	@Success		200		{object}	http.BaseResponse
//	@Failure		400		{object}	http.BaseResponse
//	@Failure		500		{object}	http.BaseResponse
//	@Router			/register [post]
func (a *AuthHttpHandler) Regis(c *gin.Context) {
	// init request body
	var regisRequest common.RegisterRequest

	//bind request body
	if err := c.ShouldBindJSON(&regisRequest); err != nil {
		httputil.WriteBadRequestResponseWithErrMsg(c, httputil.ResponseBadRequestError, err)
		return

	}

	// validate request body
	if err := validator.New().Struct(&regisRequest); err != nil {
		httputil.WriteBadRequestResponseWithErrMsg(c, httputil.ResponseBadRequestError, err)
		return
	}

	user := domain.User{
		Email:    regisRequest.Email,
		Password: regisRequest.Password,
		Username: regisRequest.Username,
		Age:      regisRequest.Age,
	}

	err := a.authUseCase.Register(user)

	if err != nil {
		httputil.WriteServerErrorResponse(c, httputil.ResponseServerError, err)
		return
	}

	httputil.WriteOkResponse(c, "Register success")
	return
}
