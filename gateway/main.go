package main

import (
	"context"
	"log"
	"net/http"
	proto "proto/proto/auth/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var openAPISpecPaths = []string{
	"openapi/api.swagger.json",
	"gen/openapiv2/api.swagger.json",
	"../gen/openapiv2/api.swagger.json",
}

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Task Engine API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`

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
	apiMux := runtime.NewServeMux()

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

	return proto.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		config.AuthHost+":"+config.AuthPort,
		opts,
	)
}
