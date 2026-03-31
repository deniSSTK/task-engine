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

	config := NewConfig()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := proto.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:"+config.AuthPort,
		opts,
	); err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", config.AppPort)
	if err := http.ListenAndServe(
		":"+config.AppPort,
		mux,
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}
