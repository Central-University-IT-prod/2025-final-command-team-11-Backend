package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"REDACTED/team-11/backend/booking/pkg/logger"
	"REDACTED/team-11/backend/booking/pkg/middlewares"
	api "REDACTED/team-11/backend/booking/pkg/ogen"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
}

func NewServer(
	ogenHandler api.Handler,
	securityHandler api.SecurityHandler,
	l *zap.Logger,
) (*Server, error) {
	ogenServer, err := api.NewServer(ogenHandler, securityHandler, api.WithErrorHandler(errorHandler))
	if err != nil {
		return nil, err
	}

	handler := middlewares.Apply(
		ogenServer,
		middlewares.Recover(),
		middlewares.LoggerProvider(l),
		middlewares.Logging(),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOD"))
	})
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", handler))

	return &Server{
		server: &http.Server{
			Handler: mux,
		},
	}, nil
}

func (s *Server) Run(ctx context.Context, port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	if err != nil {
		logger.FromCtx(ctx).Debug("handling error", zap.Error(err))
	}
	switch code {
	case http.StatusBadRequest:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(code)
}
