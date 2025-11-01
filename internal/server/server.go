package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-native-webserver/config"
	"go-native-webserver/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

var (
	instance     *Server
	instanceOnce = false
)

// GetServerInstance returns a singleton Server instance using config.
func GetServerInstance(config *config.ServerConfig) *Server {
	if instance == nil {
		port := 8080 // default
		if config != nil {
			port = configPort(config)
		}
		instance = &Server{
			port: port,
			db:   database.New(),
		}
		instanceOnce = true
	}
	return instance
}

// Helper to get port from config or env
func configPort(config *config.ServerConfig) int {
	if config != nil {
		if p := os.Getenv("PORT"); p != "" {
			if port, err := strconv.Atoi(p); err == nil {
				return port
			}
		}
	}
	return 8080
}

// NewServer creates a new http.Server using config
func NewServer(config *config.ServerConfig) *http.Server {
	srv := GetServerInstance(config)
	mux := http.NewServeMux()
	mainHandler := srv.RegisterRoutes()
	mux.Handle("/", mainHandler)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
