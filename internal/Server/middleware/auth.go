package middleware

import "net/http"

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		//Получение userId и перадача по контексту

		//
		next.ServeHTTP(w, r)

	})
}
