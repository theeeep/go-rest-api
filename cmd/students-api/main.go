package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theeeep/go-rest-api/internal/config"
)

func main() {
	// Load Config
	cfg := config.MustLoad()

	// db setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

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
