package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/trb1maker/self-sign-cert-example/internal/http/client"
)

//nolint:gochecknoglobals
var (
	certPath string
	wait     int
)

func main() {
	flag.StringVar(&certPath, "cert", "pravo.tech.pem", "Path to the certificate file")
	flag.IntVar(&wait, "wait", 5, "Time to wait before shutting down the server") //nolint:mnd
	flag.Parse()

	c, err := client.NewClient(10_000, certPath) //nolint:mnd,varnamelen
	if err != nil {
		slog.LogAttrs(
			context.Background(), slog.LevelError, "Client",
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(wait)*time.Second)
	defer cancel()

	ts, err := c.Time(ctx) //nolint:varnamelen
	if err != nil {
		slog.LogAttrs(
			context.Background(), slog.LevelError, "Client",
			slog.String("err", err.Error()),
		)
		os.Exit(1) //nolint:gocritic
	}

	slog.LogAttrs(
		context.Background(), slog.LevelInfo, "Client",
		slog.String("time", ts.String()),
	)
}
