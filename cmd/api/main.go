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

	"go-native-webserver/config"
	"go-native-webserver/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

// fml need to rewrite using cobra... messy af
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "reload-config":
			reloadConfig()
			return
		case "start-server":
			startServer()
			return
		case "start-queue-worker":
			// Implement queue worker start logic here
			log.Println("Starting queue worker...")
			return
		case "queue": // process a specific job (retry failed job)
			jobIDStr := os.Args[2]
			if jobIDStr == "" {
				log.Fatalf("Job ID must be provided for 'queue' command")
			}
			// Implement queue job processing logic here
			log.Printf("Processing queue job with ID: %s", jobIDStr)

			return
		default:
			log.Fatalf("Unknown command: %s", os.Args[1])
		}
	}
}

func startServer() {
	config := config.GetServerConfig()

	srv := server.NewServer(config)
	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	<-done
	log.Println("Graceful shutdown complete.")
}

func reloadConfig() {
	// Implement configuration reload logic here
	log.Println("Reloading configuration...")
	config.ReloadServerConfig()
}
