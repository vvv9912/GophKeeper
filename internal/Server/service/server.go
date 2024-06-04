package service

import (
	"GophKeeper/pkg/logger"
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
func StartServer(ctx context.Context, h http.Handler, addr, cert, key string) *Server {

	server := &http.Server{
		Addr:    addr,
		Handler: h,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	go func() {
		logger.Log.Info("Start server...", zap.String("addr", addr))
		err := server.ListenAndServeTLS(cert, key)
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
