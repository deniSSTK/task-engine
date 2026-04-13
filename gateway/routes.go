package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func registerRoutes(mux *http.ServeMux, apiMux *runtime.ServeMux) {
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	mux.HandleFunc("/openapi.json", openAPISpecHandler)
	mux.HandleFunc("/docs", redirectToDocs)
	mux.HandleFunc("/docs/", swaggerUIHandler)
}

func openAPISpecHandler(w http.ResponseWriter, _ *http.Request) {
	spec, err := readOpenAPISpec()
	if err != nil {
		http.Error(w, "swagger not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(spec)
}

func redirectToDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
}

func swaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func readOpenAPISpec() ([]byte, error) {
	for _, path := range openAPISpecPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
	}

	return nil, errors.New("openapi spec not found")
}
