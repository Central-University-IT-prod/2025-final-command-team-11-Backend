package v1

import (
	"github.com/gin-gonic/gin"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	"github.com/nikitaSstepanov/tools/ctx"
)

type AccountHandler interface {
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Update(ctx *gin.Context)
	SetRole(c *gin.Context)
	Delete(c *gin.Context)
	GetByEmail(c *gin.Context)
	Edit(c *gin.Context)
	GetList(c *gin.Context)
}

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

type Middleware interface {
	CheckAccess(roles ...types.Role) gin.HandlerFunc
	InitLogger(c ctx.Context) gin.HandlerFunc
}
