package handler

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

type ctxKeyRequestID struct{}

// RequestID adds or reads X-Request-ID header and puts it in context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID{}, id)))
	})
}

// Logging logs request start, end, and duration.
func Logging(logger svclog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Request(r.Context(), start.Unix(), time.Now().Unix(), r.Method+" "+r.URL.Path)
		})
	}
}

// Recovery catches panics and returns 500.
func Recovery(logger svclog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(r.Context(), "panic", fmt.Sprint(rec), string(debug.Stack()))
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal_error","message":"An error occurred"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
