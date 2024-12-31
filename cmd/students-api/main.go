package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theeeep/go-rest-api/internal/config"
	"github.com/theeeep/go-rest-api/internal/http/handlers/student"
	"github.com/theeeep/go-rest-api/internal/storage/sqlite"
)

func main() {
	// Load Config
	cfg := config.MustLoad()

	// setup storage / db
	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/students", student.New())

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Starting server", "addr", fmt.Sprintf("http://%s", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			slog.Error("Failed to start server: ", slog.String("error", err.Error()))
		}
	}() // () --> anonymous function means it will be executed immediately

	<-done

	slog.Info("Shutting down server...") // slog --> structured logging

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // context(ctx) --> allows us to cancel operations
	defer cancel()                                                          // defer --> executes a function at the end

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server: ", slog.String("error", err.Error()))
	}

	slog.Info("Server exited properly")
}
