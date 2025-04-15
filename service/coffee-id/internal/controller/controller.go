package controller

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/httper"
)

type Config struct {
	V1   v1.Config `yaml:"v1"`
	Mode string    `yaml:"mode" env:"MODE" env-default:"DEBUG"`
}

type Controller struct {
	v1  *v1.Router
	cfg *Config
}

func New(uc *usecase.UseCase, cfg *Config) *Controller {

	return &Controller{
		v1:  v1.New(uc, &cfg.V1),
		cfg: cfg,
	}
}

func (c *Controller) InitRoutes(ctx ctx.Context) *gin.Engine {
	setGinMode(c.cfg.Mode)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func (origin string) bool {return true},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Yandex-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(httper.StatusOK, "pong")
	})

	api := router.Group("/api")
	{
		c.v1.InitRoutes(ctx, api)
	}

	return router
}

func setGinMode(mode string) {
	switch mode {

	case "RELEASE":
		gin.SetMode(gin.ReleaseMode)

	case "TEST":
		gin.SetMode(gin.TestMode)

	case "DEBUG":
		gin.SetMode(gin.DebugMode)

	default:
		gin.SetMode(gin.DebugMode)

	}
}
