package service

import (
	"GophKeeper/pkg/logger"
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"net/http"
)

// StartServer - запуск сервера.
func StartServer(ctx context.Context, h http.Handler, addr, cert, key string) *http.Server {

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
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Log.Error("shutdown error", zap.Error(err))
		}
	}()
	return server
}
