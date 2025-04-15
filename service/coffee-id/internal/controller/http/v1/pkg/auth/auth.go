package auth

import (
	"github.com/gin-gonic/gin"
	conv "github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/converter"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/dto"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/validator"
	resp "github.com/nikitaSstepanov/coffee-id/internal/controller/response"
	ct "github.com/nikitaSstepanov/coffee-id/pkg/utils/controller"
	"github.com/nikitaSstepanov/tools/httper"
)

type Auth struct {
	usecase AuthUseCase
	cookie  *httper.CookieCfg
	host    string
}

func New(uc AuthUseCase, cookie *httper.CookieCfg, host string) *Auth {
	return &Auth{
		usecase: uc,
		cookie:  cookie,
		host:    host,
	}
}

// @Summary Log in a user
// @Description Logs in a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.Login true "Login information"
// @Success 200 {object} dto.Token "Access token"
// @Failure 400 {object} resp.JsonError "Incorrect data"
// @Failure 401 {object} resp.JsonError "Incorrect email or password"
// @Failure 404 {object} resp.JsonError "This user wasn't found."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body, validator.Password); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user := conv.EntityLogin(body)

	tokens, user, err := a.usecase.Login(ctx, user)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.SetCookie(
		a.cookie.Name, tokens.Refresh, a.cookie.Age,
		a.cookie.Path, a.cookie.Host,
		a.cookie.Secure, a.cookie.HttpOnly,
	)

	result := conv.DtoAnswer(user, tokens.Access)

	c.JSON(okStatus, result)
}

// @Summary Log out a user
// @Description Logs out a user by invalidating the session
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} resp.Message "Logout success."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router  /id/auth/logout [post]
func (a *Auth) Logout(c *gin.Context) {
	c.SetCookie(
		a.cookie.Name, "", -1,
		a.cookie.Path, a.cookie.Host,
		a.cookie.Secure, a.cookie.HttpOnly,
	)

	c.JSON(okStatus, logoutMsg)
}

// @Summary Refresh user tokens
// @Description Refreshes the user's tokens using the refresh token from the cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.Token "Refresh token"
// @Failure 401 {object} resp.JsonError "Token is invalid"
// @Failure 404 {object} resp.JsonError "Your token wasn't found., This user wasn't found."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/auth/refresh [get]
func (a *Auth) Refresh(c *gin.Context) {
	ctx := ct.GetCtx(c)

	refresh, cookieErr := c.Cookie(a.cookie.Name)
	if cookieErr != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(cookieErr))
		return
	}

	tokens, err := a.usecase.Refresh(ctx, refresh)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.SetCookie(
		a.cookie.Name, tokens.Refresh, a.cookie.Age,
		a.cookie.Path, a.cookie.Host,
		a.cookie.Secure, a.cookie.HttpOnly,
	)

	result := conv.DtoToken(tokens)

	c.JSON(okStatus, result)
}
