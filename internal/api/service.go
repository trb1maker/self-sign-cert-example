//go:generate protoc --proto_path=../../api --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ../../api/api.proto

package api

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"github.com/trb1maker/self-sign-cert-example/internal/http/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func StartService(ctx context.Context, certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to load key pair: %w", err)
	}

	lis, err := net.Listen("tcp", ":10000") //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	server := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ensureValidBasicCredentials),
	)
	RegisterServiceServer(server, &Server{})

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

type Server struct {
	UnimplementedServiceServer
}

func (s *Server) Hello(_ context.Context, _ *Request) (*Response, error) {
	return &Response{Message: "Hello!"}, nil
}

func valid(authorization []string) bool {
	if len(authorization) == 0 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Basic ")

	return token == base64.StdEncoding.EncodeToString([]byte(creds.User+":"+creds.Password))
}

func ensureValidBasicCredentials(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}

	return handler(ctx, req)
}
