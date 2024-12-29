package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/self-sign-cert-example/internal/http/creds"
)

func TestAuth(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		authMiddleware(testHandler).ServeHTTP(w, r)
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("authorized", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.SetBasicAuth(creds.User, creds.Password)

		authMiddleware(testHandler).ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
}
