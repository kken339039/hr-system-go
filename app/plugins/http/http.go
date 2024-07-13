package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewHttpServer)
}

type HttpServer struct {
	logger *logger.Logger
	env    *env.Env
	srv    *http.Server
}

func NewHttpServer(logger *logger.Logger, env *env.Env, lc fx.Lifecycle, engine *gin.Engine) *HttpServer {
	port := env.GetEnv("PORT")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: engine,
	}

	server := &HttpServer{
		logger: logger,
		env:    env,
		srv:    srv,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			server.Serve()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown(ctx)
			return nil
		},
	})
	return server
}

func (s HttpServer) Serve() {
	s.logger.Info("Starting HTTP server")
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()
	s.logger.Info("HTTP Server is up and running", zap.String("addr", s.srv.Addr))
}

func (s HttpServer) Shutdown(ctx context.Context) {
	s.logger.Info("Shutting down HTTP server")
	_ = s.srv.Shutdown(ctx)
}
