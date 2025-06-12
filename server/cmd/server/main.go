package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devmuggs/onigiri/server/internal/config"
	"github.com/devmuggs/onigiri/server/internal/db"
	"github.com/devmuggs/onigiri/server/internal/features/auth"
	"github.com/devmuggs/onigiri/server/internal/features/users"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func run() error {

	var err error

	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to init zap: %v", err)
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := config.New(ctx, logger)
	if err != nil {
		logger.Fatal("failed to load config:", zap.Error(err))
	}

	database, err := db.New(ctx, config.DBConnectionString)
	if err != nil {
		logger.Fatal("failed to initialise db:", zap.Error(err))
	}
	defer database.Close()

	r := chi.NewRouter()
	r.Get("/api/hello-world", helloHandler)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", users.NewRouter(database.Pool, logger))
		r.Mount("/auth", auth.NewRouter(database.Pool, logger))
	})

	logger.Info("Starting Onigiri server")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting Onigiri server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}()

	<-stop
	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited cleanly")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling /hello-world", zap.String("method", r.Method), zap.String("remote", r.RemoteAddr))
	fmt.Fprint(w, "Hello World")
}
