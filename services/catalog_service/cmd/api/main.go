package main

import (
	"catalog-service/internal/database"
	"catalog-service/internal/event"
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
		Host:       getEnv("DB_HOST", "localhost"),
		Port:       getEnv("DB_PORT", "5432"),
		User:       getEnv("DB_USER", "postgres"),
		Password:   getEnv("DB_PASSWORD", "123456"),
		DBName:     getEnv("DB_NAME", "postgres"),
		SearchPath: getEnv("DB_SEARCH_PATH", "public"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresPool(dbCfg)
	if err != nil {
		slog.Error("failed to connect to postgres: ", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("connected to postgres")

	eventRepo := event.NewRepository(db)
	eventHandler := event.NewHandler(eventRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /events", eventHandler.GetAllHandler)
	mux.HandleFunc("GET /events/{id}", eventHandler.GetEventHandler)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
		os.Exit(1)
	}
	slog.Info("server shutdown successfully")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
