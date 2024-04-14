package main

import (
	"banners/internal/config"
	hand "banners/internal/handler"
	serv "banners/internal/service"
	"banners/internal/storage/postgresql"
	"banners/internal/storage/redisC"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting banners server",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)

	log.Debug("debug messages are enabled")

	repo, err := postgresql.New(cfg.DataSourceName, 5)
	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}

	c, err := redisC.New(cfg.CacheStorage.Address)
	if err != nil {
		log.Error("failed to initialize cache", err)
		os.Exit(1)
	}

	service, err := serv.New(log, repo, repo, c)
	if err != nil {
		log.Error("failed to initialize services", err)
		os.Exit(1)
	}

	deleteCtx, _ := context.WithCancel(context.Background())
	handler, err := hand.New(log, service, service, service, deleteCtx)
	if err != nil {
		log.Error("failed to initialize handlers", err)
		os.Exit(1)
	}

	router := handler.InitRoutes()

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:        ":" + cfg.HTTPServer.Address,
		Handler:     router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.Timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)
		return
	}

	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
