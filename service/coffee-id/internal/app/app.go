package app

import (
	"fmt"

	"github.com/nikitaSstepanov/coffee-id/internal/controller"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/storage"
	"github.com/nikitaSstepanov/tools"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	"github.com/nikitaSstepanov/tools/sl"
)

type App struct {
	controller *controller.Controller
	usecase    *usecase.UseCase
	storage    *storage.Storage
	server     *httper.Server
	ctx        ctx.Context
}

func New() *App {
	if err := tools.Init(); err != nil {
		panic(fmt.Sprintf("Can`t init tools. Error: %v", err))
	}

	log := tools.Sl()

	cfg, err := getConfig()
	if err != nil {
		log.Error("Can`t get app config", sl.ErrAttr(err))
	}

	ctx := ctx.New(log)

	app := &App{}

	app.ctx = ctx

	app.storage = storage.New(ctx)

	app.usecase = usecase.New(app.storage, &cfg.UseCase)

	app.controller = controller.New(app.usecase, &cfg.Controller)

	handler := app.controller.InitRoutes(ctx)

	app.server = tools.HttpServer(handler)

	return app
}

func (a *App) Run() {
	log := a.ctx.Logger()

	a.server.Start()

	if err := a.shutdown(); err != nil {
		log.Error("Failed to shutdown app", sl.ErrAttr(err))
		return
	}

	log.Info("Application stopped successfully")
}

func (a *App) shutdown() e.Error {
	log := a.ctx.Logger()

	err := e.E(a.server.Shutdown(a.ctx))
	if err != nil {
		log.Error("Failed to stop http server", err.SlErr())
		return err
	}

	return nil
}
