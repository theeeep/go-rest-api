package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Printf("🌐 Server started on %s\n", cfg.HTTPServer.Addr)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}() // () --> anonymous function means it will be executed immediately

}
