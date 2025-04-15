package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

type EntityHandler interface {
	Save(c *gin.Context)
	GetEntities(c *gin.Context)
	GetFloors(c *gin.Context)
	DeleteFloor(c *gin.Context)
	EntityById(c *gin.Context)
}

type GuestHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
}

type BookingHandler interface {
	CheckAccess(c *gin.Context)
	Stats(c *gin.Context)
}

type OrderHandler interface {
	SetStatus(c *gin.Context)
	Get(c *gin.Context)
	Stats(c *gin.Context)
}

type VerificationHandler interface {
	CheckVerify(c *gin.Context)
	Verify(c *gin.Context)
}

type Middleware interface {
	CheckAccess(roles ...types.Role) gin.HandlerFunc
	InitLogger(c ctx.Context) gin.HandlerFunc
}
