package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"REDACTED/team-11/backend/booking/internal/config"
	coffeeid "REDACTED/team-11/backend/booking/internal/repo/coffee-id"
	"REDACTED/team-11/backend/booking/internal/repo/postgres"
	"REDACTED/team-11/backend/booking/internal/service"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1/handlers"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1/security"
	"REDACTED/team-11/backend/booking/pkg/logger"
	pg_helper "REDACTED/team-11/backend/booking/pkg/postgres"
	"go.uber.org/zap"
)

var (
	shutdownTimeout = time.Second * 5
)

func main() {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal("get config:", err)
	}

	l, err := logger.Get(cfg.LogLevel)
	if err != nil {
		log.Fatal("get logger:", err)
	}

	db, err := pg_helper.Connect(ctx, cfg.PostgresConfig)
	if err != nil {
		l.Fatal("connect to postgres", zap.Error(err))
	}

	floorsRepo := postgres.NewFloorsRepo(db)
	usersRepo := coffeeid.NewUserRepo(cfg.CoffeeIdBaseUrl)
	bookingEntitiesRepo := postgres.NewBookingEntitiesRepo(db)
	bookingsRepo := postgres.NewBookingsRepo(db)
	ordersRepo := postgres.NewOrdersRepo(db)

	workloadsService := service.NewWorkloadService(bookingEntitiesRepo, bookingsRepo, floorsRepo)
	ordersService := service.NewOrdersService(ordersRepo, bookingsRepo)
	bookingsService := service.NewBookingsService(bookingsRepo, bookingEntitiesRepo, ordersRepo, workloadsService, usersRepo)

	bookingsHandler := handlers.NewBookingsHandler(bookingsService)
	ordersHandler := handlers.NewOrdersHandler(ordersService)
	workloadsHandler := handlers.NewWorkloadsHandler(workloadsService)

	securityHandler := security.NewSecurityHandler(cfg.JWTSecret)
	handler := http.NewHandler(
		bookingsHandler,
		ordersHandler,
		workloadsHandler,
	)

	server, err := http.NewServer(handler, securityHandler, l)
	if err != nil {
		l.Fatal("get server", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		l.Info("starting http server", zap.Int("port", cfg.ServerPort))
		if err := server.Run(ctx, cfg.ServerPort); err != nil && err != stdhttp.ErrServerClosed {
			l.Fatal("start http server", zap.Error(err))
		}
	}()

	<-sigCh

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	l.Info("shutting down http server")
	if err := server.Shutdown(ctx); err != nil {
		l.Error("shutdown http server", zap.Error(err))
	}

	l.Info("server stopped")
}
