package service

import (
	"GophKeeper/pkg/logger"
	"context"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func StartServer(ctx context.Context, h http.Handler) *Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Log.Error("Start server error", zap.Error(err))
			return
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			err := server.Shutdown(ctx)
			if err != nil {
				logger.Log.Error("shutdown error", zap.Error(err))
				return
			}
		}

	}()
	return &Server{}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
