package api

import (
	context "context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/trb1maker/self-sign-cert-example/internal/http/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type basicAuth struct {
	username string
	password string
}

func (b basicAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(b.username+":"+b.password)),
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}

func StartClient(ctx context.Context, certPath string) error {
	cert, err := credentials.NewClientTLSFromFile(certPath, "localhost")
	if err != nil {
		return fmt.Errorf("failed to read cert file: %w", err)
	}

	auth := basicAuth{
		username: creds.User,
		password: creds.Password,
	}

	conn, err := grpc.NewClient(
		"localhost:10000",
		grpc.WithTransportCredentials(cert),
		grpc.WithPerRPCCredentials(auth),
	)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}
	defer conn.Close()

	client := NewServiceClient(conn)

	resp, err := client.Hello(ctx, &Request{})
	if err != nil {
		return fmt.Errorf("failed to call Hello: %w", err)
	}

	//nolint:contextcheck
	slog.LogAttrs(
		context.Background(), slog.LevelInfo, "Client",
		slog.String("message", resp.GetMessage()),
	)

	return nil
}
