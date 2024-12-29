package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/trb1maker/self-sign-cert-example/internal/api"
)

//nolint:gochecknoglobals
var (
	certPath, keyPath string
	wait              int
)

func main() {
	flag.StringVar(&certPath, "cert", "pravo.tech.pem", "Path to the certificate file")
	flag.StringVar(&keyPath, "key", "pravo.tech.key", "Path to the key file")
	flag.IntVar(&wait, "wait", 5, "Time to wait before shutting down the server") //nolint:mnd
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(wait))
	defer cancel()

	if err := api.StartService(ctx, certPath, keyPath); err != nil {
		slog.LogAttrs(
			context.Background(), slog.LevelError, "Server",
			slog.String("err", err.Error()),
		)

		os.Exit(1) //nolint:gocritic
	}
}