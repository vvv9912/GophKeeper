package service

import (
	"GophKeeper/pkg/logger"
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func StartServer(ctx context.Context) {
	r := chi.NewRouter()
	//http.Serve(autocert.NewListener())
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
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
}
