package account

import (
	"strconv"

	"github.com/gin-gonic/gin"
	conv "github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/converter"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/dto"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/validator"
	resp "github.com/nikitaSstepanov/coffee-id/internal/controller/response"
	ct "github.com/nikitaSstepanov/coffee-id/pkg/utils/controller"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
)

type Account struct {
	usecase AccountUseCase
	cookie  *httper.CookieCfg
}

func New(uc AccountUseCase, cookie *httper.CookieCfg) *Account {
	return &Account{
		usecase: uc,
		cookie:  cookie,
	}
}

// @Summary Retrieve user own account
// @Description Returns user own account.
// @Tags Account
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.Account "Successful response"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 404 {object} resp.JsonError "This user wasn`t found."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/ [get]
func (a *Account) Get(c *gin.Context) {
	ctx := ct.GetCtx(c)

	userId := c.GetString("userId")

	user, err := a.usecase.Get(ctx, userId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoUser(user)

	c.JSON(okStatus, result)
}

// @Summary Retrieve user by ID
// @Description Returns user information based on their ID.
// @Tags Account
// @Accept json
// @Produce json
// @Security Bearer
// @Param        id    path     string  true  "user id"  Format(uuid)
// @Success 200 {object} dto.Account "Successful response"
// @Failure 400 {object} resp.JsonError "ID must be integer"
// @Failure 404 {object} resp.JsonError "This user wasn`t found."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/{id} [get]
func (a *Account) GetById(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user, err := a.usecase.Get(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoUser(user)

	c.JSON(okStatus, result)
}

// @Summary Retrieve user by Email
// @Description Returns user information based on their Email.
// @Tags Account
// @Accept json
// @Produce json
// @Security Bearer
// @Param        email    path     string  true  "user email"  Format(email)
// @Success 200 {object} dto.Account "Successful response"
// @Failure 400 {object} resp.JsonError "Invalid email"
// @Failure 404 {object} resp.JsonError "This user wasn`t found."
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/email/{email} [get]
func (a *Account) GetByEmail(c *gin.Context) {
	ctx := ct.GetCtx(c)

	email := c.Param("email")

	if err := validator.Email(email); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user, err := a.usecase.GetByEmail(ctx, email)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoUser(user)

	c.JSON(httper.StatusOK, result)
}

// @Summary Create User
// @Description Creates a new user and returns access tokens.
// @Tags Account
// @Accept json
// @Produce json
// @Param body body dto.CreateUser true "Data for creating a user"
// @Success 201 {object} dto.AccountAnswer "Successful response with token"
// @Failure 400 {object} resp.JsonError "Incorrect data"
// @Failure 409 {object} resp.JsonError "User with this email already exist"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/new [post]
func (a *Account) Create(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.CreateUser

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body, validator.Password, validator.Birthday); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user := conv.EntityCreate(body)

	tokens, err := a.usecase.Create(ctx, user)
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

	c.JSON(createdStatus, result)
}

// @Summary Update user information
// @Description Updates the user's information including password.
// @Tags Account
// @Accept json
// @Produce json
// @Param body body dto.UpdateUser true "User update data"
// @Security Bearer
// @Success 200 {object} dto.Account "Updated."
// @Failure 400 {object} resp.JsonError "Incorrect data."
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 404 {object} resp.JsonError "This user wasn't found"
// @Failure 409 {object} resp.JsonError "User with this email already exists"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/edit [patch]
func (a *Account) Update(c *gin.Context) {
	ctx := ct.GetCtx(c)

	userId := c.GetString("userId")

	var body dto.UpdateUser

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body, validator.Password, validator.Birthday); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user := conv.EntityUpdate(body)
	user.Id = userId

	user, err := a.usecase.Update(ctx, user, body.OldPassword)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoUser(user)

	c.JSON(okStatus, result)
}

func (a *Account) SetRole(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.SetRole

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	user := conv.EntitySetRole(body)

	err := a.usecase.AddRole(ctx, user)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(okStatus, okMsg)
}

// @Summary Update user information
// @Description Updates the user's information including password. For admin
// @Tags Account
// @Accept json
// @Produce json
// @Param body body dto.UpdateUser true "User update data"
// @Security Bearer
// @Param        id    path     string  true  "user id"  Format(uuid)
// @Success 200 {object} dto.Account "Updated."
// @Failure 400 {object} resp.JsonError "Incorrect data."
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 404 {object} resp.JsonError "This user wasn't found"
// @Failure 409 {object} resp.JsonError "User with this email already exists"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/{id}/edit [patch]
func (a *Account) Edit(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	var body dto.UpdateUser

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body, validator.Password, validator.Birthday); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user := conv.EntityUpdate(body)
	user.Id = id

	user, err := a.usecase.Update(ctx, user, body.OldPassword)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoUser(user)

	c.JSON(okStatus, result)
}

func (a *Account) Delete(c *gin.Context) {
	ctx := ct.GetCtx(c)
	userId := c.GetString("userId")

	var body dto.DeleteUser

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(err))
		return
	}

	if err := validator.Struct(body, validator.Password); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	user := conv.EntityDelete(body)
	user.Id = userId

	err := a.usecase.Delete(ctx, user)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(noContentStatus, nil)
}

// @Summary Get list of users
// @Description List of users with pagination. Only for ADMIN
// @Tags Account
// @Produce json
// @Security Bearer
// @Param page      query int  false  "Page"
// @Param size      query int  false  "Size"
// @Success 200 {object} dto.AcountsPagintation "OK."
// @Failure 400 {object} resp.JsonError "Incorrect data."
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 404 {object} resp.JsonError "This user wasn't found"
// @Failure 409 {object} resp.JsonError "User with this email already exists"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /id/account/all [get]
func (a *Account) GetList(c *gin.Context) {
	ctx := ct.GetCtx(c)

	token := c.GetHeader("Authorization")

	page, parseErr := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 64)
	if parseErr != nil || page < 0 {
		resp.AbortErrMsg(c, e.New("Page must be integer", e.BadInput))
		return
	}

	size, parseErr := strconv.ParseInt(c.DefaultQuery("size", "5"), 10, 64)
	if parseErr != nil || size < 0 {
		resp.AbortErrMsg(c, e.New("Page must be integer", e.BadInput))
		return
	}

	list, count, err := a.usecase.GetList(ctx, int(page), int(size), token)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := &dto.AcountsPagintation{
		Accounts: list,
		Count:    count,
	}

	c.JSON(okStatus, result)
}
