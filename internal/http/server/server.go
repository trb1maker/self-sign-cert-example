package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func NewServer(port uint16, certPath, keyPath string) *Sever {
	mux := http.NewServeMux()

	mux.Handle("GET /now", authMiddleware(http.HandlerFunc(timeHandler)))

	return &Sever{
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			ReadHeaderTimeout: time.Second * 5, //nolint:mnd
			Handler:           mux,
		},
		certPath: certPath,
		keyPath:  keyPath,
	}
}

type Sever struct {
	srv               *http.Server
	certPath, keyPath string
}

func (s *Sever) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //nolint:mnd
		defer cancel()

		//nolint:contextcheck
		if err := s.srv.Shutdown(ctx); err != nil {
			slog.LogAttrs(
				context.Background(), slog.LevelError, "can't shutdown server",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	if err := s.srv.ListenAndServeTLS(s.certPath, s.keyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("can't start server: %w", err)
	}

	return nil
}
