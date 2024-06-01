package middleware

import (
	"GophKeeper/internal/Server/service"
	"context"
	"net/http"
	"strings"
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

		// Отделение префикса "Bearer" от токена
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}
		userId, err := m.GetUserId(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userId == -1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//todo
		ctx := context.WithValue(r.Context(), "UserId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
