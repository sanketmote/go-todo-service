package main

import (
	"context"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/go-todo-service/internal/config"
	"github.com/sanketmote/go-todo-service/internal/handler"
	"github.com/sanketmote/go-todo-service/internal/repository"
	"github.com/sanketmote/go-todo-service/internal/todo"
	"github.com/sanketmote/go-todo-service/internal/upgrade"
	"github.com/sanketmote/gokit-wrapper/httptransport"
	"github.com/sanketmote/gokit-wrapper/httptransport/healthcheck"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		os.Exit(1)
	}

	logger := svclog.NewStdoutCtxLogger("go-todo-service")
	ctx := context.Background()

	db, err := sqlx.Connect("mysql", cfg.DSN)
	if err != nil {
		logger.Error(ctx, "db_connect", err.Error(), "")
		os.Exit(1)
	}
	defer db.Close()

	if err := upgrade.Run(ctx, db, logger); err != nil {
		logger.Error(ctx, "upgrade", err.Error(), "")
		os.Exit(1)
	}

	repo := repository.NewTodoRepository(db)
	svc := todo.NewTodoService(repo, logger)
	endpoints := todo.MakeEndpoints(svc)

	r := mux.NewRouter()

	// Middleware: recovery (outer) → logging → request id → routes
	r.Use(handler.Recovery(logger))
	r.Use(handler.Logging(logger))
	r.Use(handler.RequestID)

	// Health
	healthToken := cfg.HealthToken
	if healthToken == "" {
		healthToken = "dev"
	}
	version := cfg.Version
	if version == "" {
		version = "1.0.0"
	}
	if err := healthcheck.AddHealthCheckAPIHandlers(r, healthcheck.HealthConfig{
		Logger:           logger,
		ServiceName:      "go-todo-service",
		XdrvHealthzToken: healthToken,
		Version:          version,
		BuildInfo:        cfg.BuildInfo,
		StartTime:        cfg.StartTime,
		MIMETYPE:         httptransport.MIMEJSON,
	}); err != nil {
		logger.Error(ctx, "health_setup", err.Error(), "")
		os.Exit(1)
	}

	// Todo API: path → decode → endpoint → encode
	todo.AddHandlers(r, endpoints)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	logger.Info(ctx, "layer", "main", "msg", "Starting application server")

	err = server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info(ctx, "server", "Server closed gracefully")
		} else {
			logger.Error(ctx, "server", err.Error(), "")
			os.Exit(1)
		}
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "shutdown", err.Error(), "")
	}
	logger.Info(ctx, "server", "Server shutdown gracefully")
}
