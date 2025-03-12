package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LincolnG4/Haku/internal/controllers"
	"github.com/LincolnG4/Haku/internal/services"
)

func main() {
	// Configuration (could be from env vars, config file, etc.)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Services
	pipelineService := services.NewPipelineService()
	defer pipelineService.Close()

	// Controllers
	pipeline := controllers.NewPipelineController(pipelineService)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pipelines", pipeline.CreatePipeline)

	// Server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Start the server
	log.Printf("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped.")
}
