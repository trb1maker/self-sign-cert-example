package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/trb1maker/self-sign-cert-example/internal/cert"
)

type certificateCreator interface {
	CreateCertificate(dirName string, domens ...string) error
}

func main() {
	var certPath string

	flag.StringVar(&certPath, "out", ".", "certificate path")
	flag.Parse()

	var generator certificateCreator = new(cert.Generator)

	if err := generator.CreateCertificate(certPath); err != nil {
		slog.LogAttrs(
			context.Background(), slog.LevelError, "Generator",
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}
}
