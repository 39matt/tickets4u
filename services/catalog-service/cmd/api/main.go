package main

import (
	"catalog-service/internal/database"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbCfg := database.Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "123456",
		DBName:   "events",
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresPool(dbCfg)
	if err != nil {
		slog.Error("failed to connect to postgres: ", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("connected to postgres")

	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", port)
		if err = server.ListenAndServe(); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}
	slog.Info("server shutdown successfully")
}
