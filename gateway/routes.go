package main

import (
	_ "embed"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

//go:embed openapi/api.swagger.json
var openAPISpec []byte

//go:embed swagger/index.html
var swaggerUIHTML []byte

func registerRoutes(mux *http.ServeMux, apiMux *runtime.ServeMux) {
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	mux.HandleFunc("/api/openapi.json", openAPISpecHandler)
	mux.HandleFunc("/api/docs", redirectToDocs)
	mux.HandleFunc("/api/docs/", swaggerUIHandler)
}

func openAPISpecHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(openAPISpec)
}

func redirectToDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/docs/", http.StatusMovedPermanently)
}

func swaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(swaggerUIHTML)
}
