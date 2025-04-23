package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimasbagussusilo/go-clean-boilerplate/config"
	httpDelivery "github.com/dimasbagussusilo/go-clean-boilerplate/delivery/http"
	"github.com/dimasbagussusilo/go-clean-boilerplate/delivery/http/middleware"
	"github.com/dimasbagussusilo/go-clean-boilerplate/infrastructure/repository/memory"
	"github.com/dimasbagussusilo/go-clean-boilerplate/usecase"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)
	logger.Println("Starting server...")

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize repositories
	userRepo := memory.NewUserRepository()

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize HTTP handlers
	userHandler := httpDelivery.NewUserHandler(userUseCase)

	// Create router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	userHandler.RegisterRoutes(mux)

	// Apply middleware
	handler := middleware.Logger(logger)(mux)
	handler = middleware.ErrorHandler(logger)(handler)

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start a server in a goroutine
	go func() {
		logger.Printf("Server listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Implement proper shutdown with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown error: %v", err)
	}

	logger.Println("Server stopped")
}
