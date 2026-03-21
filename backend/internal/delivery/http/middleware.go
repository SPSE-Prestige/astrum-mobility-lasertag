package http

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware validates the Bearer token and injects the user into context.
func AuthMiddleware(authUC *usecase.AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				writeError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			user, err := authUC.ValidateToken(r.Context(), token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSMiddleware adds CORS headers for development.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs incoming requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
