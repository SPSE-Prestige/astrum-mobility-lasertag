package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/ws"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	adminHandler *AdminHandler,
	gameHandler *GameHandler,
	playerHandler *PlayerHandler,
	deviceHandler *DeviceHandler,
	authMiddleware *AuthMiddleware,
	wsHandler *ws.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("POST /api/admin/login", adminHandler.Login)
	mux.HandleFunc("POST /api/events/hit", deviceHandler.HandleHit)
	mux.HandleFunc("POST /api/games/{id}/join", playerHandler.Join)
	mux.HandleFunc("GET /api/games/{id}/leaderboard", playerHandler.Leaderboard)
	mux.HandleFunc("GET /api/games/{id}/players", playerHandler.ListPlayers)

	// WebSocket
	mux.HandleFunc("GET /ws/game/{id}", wsHandler.ServeWS)

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Admin-protected endpoints
	admin := http.NewServeMux()
	admin.HandleFunc("POST /api/games", gameHandler.Create)
	admin.HandleFunc("GET /api/games", gameHandler.List)
	admin.HandleFunc("GET /api/games/{id}", gameHandler.Get)
	admin.HandleFunc("GET /api/games/{id}/state", gameHandler.GetState)
	admin.HandleFunc("POST /api/games/{id}/start", gameHandler.Start)
	admin.HandleFunc("POST /api/games/{id}/pause", gameHandler.Pause)
	admin.HandleFunc("POST /api/games/{id}/end", gameHandler.End)
	admin.HandleFunc("POST /api/games/{id}/control", gameHandler.Control)
	admin.HandleFunc("POST /api/games/{id}/teams", gameHandler.CreateTeam)

	mux.Handle("/api/games/", authMiddleware.RequireAdmin(admin))
	mux.Handle("/api/games", authMiddleware.RequireAdmin(admin))

	// CORS wrapper
	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
