package main

import (
	"net/http"
	"time"

	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId := uuid.New().String()

		w.Header().Set("X-Request-Id", requestId)

		recorder := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}

		log.Info("incoming request",
			zap.String("request_id", requestId),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_ip", r.RemoteAddr),
		)

		next.ServeHTTP(recorder, r)

		log.Info("request completed",
			zap.String("request_id", requestId),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", recorder.Status),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
