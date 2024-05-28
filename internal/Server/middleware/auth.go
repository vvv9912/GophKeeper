package middleware

import (
	"GophKeeper/internal/Server/service"
	"context"
	"net/http"
)

type Mw struct {
	service.Auth
}

func (m Mw) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := m.GetUserId(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//todo
		ctx := context.WithValue(r.Context(), "UserId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
