package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/httper"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/middleware"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/pkg/booking"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/pkg/booking_entity"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/pkg/guest"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/pkg/order"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/pkg/verification"
	"REDACTED/team-11/backend/admin/internal/usecase"
	"REDACTED/team-11/backend/admin/pkg/swagger"
)

type Config struct {
	Cookie         httper.CookieCfg    `yaml:"cookie"`
	FrontendHost   string              `env:"FRONTEND_HOST"`
	Swagger        swagger.SwaggerSpec `yaml:"swagger"`
	SessionsSecret string              `env:"SESSIONS_SECRET"`
}

type Router struct {
	booking      BookingHandler
	entity       EntityHandler
	verification VerificationHandler
	guest        GuestHandler
	order        OrderHandler
	mid          Middleware
}

func New(uc *usecase.UseCase, cfg *Config) *Router {
	swagger.SetSwaggerConfig(cfg.Swagger)

	return &Router{
		booking:      booking.New(uc.Booking),
		entity:       booking_entity.New(uc.BookingEntity),
		guest:        guest.New(uc.Guest),
		mid:          middleware.New(uc.Auth),
		order:        order.New(uc.Order),
		verification: verification.New(uc.Verification),
	}
}

func (r *Router) InitRoutes(c ctx.Context, h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/v1")
	{
		router.Use(r.mid.InitLogger(c))

		r.initEntityRoutes(router)
		r.initVerificationRoutes(router)
		r.initOrderRoutes(router)
		booking := r.initBookingRoutes(router)
		r.initGuestsRouets(booking)
		r.initSwaggerRoute(router)
	}

	return router
}

func (r *Router) initBookingRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/booking")
	{
		router.GET("/:id/access", r.mid.CheckAccess("ADMIN"), r.booking.CheckAccess)
		router.GET("/stats", r.mid.CheckAccess("ADMIN"), r.booking.Stats)
	}

	return router
}

func (r *Router) initGuestsRouets(h *gin.RouterGroup) *gin.RouterGroup {
	h.GET("/:id/guests", r.mid.CheckAccess(), r.guest.Get)
	h.POST("/:id/guests", r.mid.CheckAccess(), r.guest.Create)
	h.DELETE("/:id/guests/:userId", r.mid.CheckAccess(), r.guest.Delete)

	return h
}

func (r *Router) initEntityRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/layout")
	{
		router.GET("/floors", r.entity.GetFloors)
		router.GET("/floors/:id", r.entity.GetEntities)
		router.POST("/floors", r.mid.CheckAccess("ADMIN"), r.entity.Save)
		router.DELETE("/floors/:id", r.mid.CheckAccess("ADMIN"), r.entity.DeleteFloor)
		router.GET("/entities/:id", r.entity.EntityById)
	}

	return router
}

func (r *Router) initVerificationRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/verification")
	{
		router.GET("/:id/check", r.mid.CheckAccess("ADMIN"), r.verification.CheckVerify)
		router.POST("/:id/set", r.mid.CheckAccess("ADMIN"), r.verification.Verify)
	}

	return router
}

func (r *Router) initOrderRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	h.GET("/orders", r.mid.CheckAccess("ADMIN"), r.order.Get)

	router := h.Group("/orders")
	{
		router.POST("/:id", r.mid.CheckAccess("ADMIN"), r.order.SetStatus)
		router.GET("/stats", r.mid.CheckAccess("ADMIN"), r.order.Stats)
	}

	return router
}

func (r *Router) initSwaggerRoute(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("swagger")
	{
		router.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return router
}
