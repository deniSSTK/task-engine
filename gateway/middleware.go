package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId := uuid.New().String()

		w.Header().Set("X-Request-Id", requestId)

		recorder := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}

		slog.Info("incoming request",
			"request_id", requestId,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", r.RemoteAddr,
		)

		next.ServeHTTP(recorder, r)

		slog.Info("request completed",
			"request_id", requestId,
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.Status,
			"duration", time.Since(start).String(),
		)
	})
}
