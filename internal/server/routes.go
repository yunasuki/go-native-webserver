package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-native-webserver/config"
	"go-native-webserver/internal/controllers"
	"go-native-webserver/internal/repositories"
	"go-native-webserver/internal/service/auth"
	"go-native-webserver/pkg/logger"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) RegisterRoutes() http.Handler {
	rootMux := http.NewServeMux()

	// public routes
	rootMux.HandleFunc("/", s.healthHandler)
	rootMux.HandleFunc("/health", s.healthHandler)

	// private routes
	apiMux := http.NewServeMux()
	// Instantiate controller
	controller := controllers.NewAllInOneController() // this line is also a possible DI...
	// ugh i am not familiar with mux package enough to do method chain routing
	apiMux.Handle("/subscriptions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			controller.PostSubscription(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		return
	}))
	apiMux.Handle("/public-holidays", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			controller.GetPublicHoliday(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}))
	apiMux.Handle("/shipping-event", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			controller.PutShippingEvent(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}))
	rootMux.Handle("/api/", http.StripPrefix("/api", s.headerValidationMiddleware(s.authMiddleware(apiMux)))) // DRY like this?

	// Wrap the mux with config reload, recovery, and CORS middleware
	return s.configReloadMiddleware(s.recoveryMiddleware(s.corsMiddleware(rootMux)))
}

// ...
// need to create new context.Context with userID value in it and pass downstream to controller function....
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 6 || authHeader[:6] != "Basic " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode base64 credentials
		encoded := authHeader[6:]
		decoded, err := auth.DecodeBasicAuth(encoded)
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		user, err := repositories.NewUserRepository().FindByEmail(decoded.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if auth.CheckPasswordHash(decoded.Password, user.PasswordHash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Add userID to context (could be username or user ID)
		ctx := context.WithValue(r.Context(), "userID", user.ID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

// configReloadMiddleware triggers hot reload of app configuration if X-Reload-Config header is present.
func (s *Server) configReloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// kinda weird to reload on every request, should have been a api or command to trigger reload the config
		config.ReloadServerConfig()

		// or not reloading but middleware for config as context?
		// ctx := context.WithValue(r.Context(), "config", config.GetServerConfig())
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// recoveryMiddleware catches panics, logs them, and returns 500 Internal Server Error.
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", zap.Any("error", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// headerValidationMiddleware checks for required HTTP headers and returns 400 if missing.
func (s *Server) headerValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := config.GetServerConfig()
		origin := r.Header.Get("Origin")
		allowed := false
		if origin == "" {
			allowed = true
		} else {
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
		}
		if !allowed {
			http.Error(w, "Invalid Origin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
