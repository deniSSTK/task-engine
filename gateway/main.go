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

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := proto.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50001",
		opts,
	); err != nil {
		panic(err)
	}

	log.Println("Central Gateway started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
