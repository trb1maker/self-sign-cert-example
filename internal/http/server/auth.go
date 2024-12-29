package server

import (
	"net/http"

	"github.com/trb1maker/self-sign-cert-example/internal/http/creds"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
		user, password, ok := r.BasicAuth()
		if !ok || user != creds.User || password != creds.Password {
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}
