package main

import (
	"context"
	"log"
	"net/http"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	config := NewConfig()
	mux, err := newGatewayMux(ctx, config)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", config.AppPort)
	log.Fatal(http.ListenAndServe(":"+config.AppPort, mux))
}

func newGatewayMux(ctx context.Context, config *Config) (*http.ServeMux, error) {
	apiMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "Authorization":
				return grpcAuth.AuthorizationHeader, true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
	)

	if err := registerAPIGateway(ctx, apiMux, config); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	registerRoutes(mux, apiMux)

	return mux, nil
}

func registerAPIGateway(ctx context.Context, mux *runtime.ServeMux, config *Config) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return authv1.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		config.AuthHost+":"+config.AuthPort,
		opts,
	)
}
