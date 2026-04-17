package main

import (
	"context"
	"net/http"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	"github.com/deniSSTK/task-engine/libs/env"
	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
	ctx := context.Background()
	config := NewConfig()
	mux, err := newGatewayMux(ctx, config)
	if err != nil {
		panic(err)
	}

	defCfg := env.NewDefConfig("GATEWAY_PORT", "../.env")

	log := logger.NewLogger(defCfg)

	handler := loggingMiddleware(mux, log)

	log.Info("starting server", zap.String("port", config.AppPort))
	if err = http.ListenAndServe(":"+config.AppPort, handler); err != nil {
		log.Fatal("failed to listen server", zap.Error(err))
	}
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
