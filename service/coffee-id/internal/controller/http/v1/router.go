package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/middleware"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/pkg/account"
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/pkg/auth"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase"
	"github.com/nikitaSstepanov/coffee-id/pkg/swagger"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/httper"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Config struct {
	Cookie         httper.CookieCfg    `yaml:"cookie"`
	FrontendHost   string              `env:"FRONTEND_HOST"`
	Swagger        swagger.SwaggerSpec `yaml:"swagger"`
	SessionsSecret string              `env:"SESSIONS_SECRET"`
}

type Router struct {
	account AccountHandler
	auth    AuthHandler
	mid     Middleware
}

func New(uc *usecase.UseCase, cfg *Config) *Router {
	swagger.SetSwaggerConfig(cfg.Swagger)

	return &Router{
		auth:    auth.New(uc.Auth, &cfg.Cookie, cfg.FrontendHost),
		account: account.New(uc.Account, &cfg.Cookie),
		mid:     middleware.New(uc.Auth),
	}
}

func (r *Router) InitRoutes(ctx ctx.Context, h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/v1")
	{
		router.Use(r.mid.InitLogger(ctx))

		r.initSwaggerRoute(router)
		r.initAccountRoutes(router)
		r.initAuthRoutes(router)
	}

	return router
}

func (r *Router) initAccountRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/account")
	{
		router.GET("/", r.mid.CheckAccess(), r.account.Get)
		router.GET("/all", r.mid.CheckAccess("ADMIN"), r.account.GetList)
		router.GET("/:id", r.account.GetById)
		router.PATCH("/:id/edit", r.mid.CheckAccess("ADMIN"), r.account.Edit)
		router.GET("/email/:email", r.account.GetByEmail)
		router.POST("/new", r.account.Create)
		router.PATCH("/edit", r.mid.CheckAccess(), r.account.Update)
		router.PATCH("/role", r.mid.CheckAccess("ADMIN"), r.account.SetRole)
		router.DELETE("/delete", r.mid.CheckAccess(), r.account.Delete)
	}

	return router
}

func (r *Router) initAuthRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/auth")
	{
		router.POST("/login", r.auth.Login)
		router.POST("/logout", r.mid.CheckAccess(), r.auth.Logout)
		router.GET("/refresh", r.auth.Refresh)
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
