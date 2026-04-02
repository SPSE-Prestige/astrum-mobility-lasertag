package http

import (
	"database/sql"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/ws"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	mqttclient "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/infrastructure/mqtt"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	cfg *config.Config,
	authUC domain.AuthUseCasePort,
	gameHandler *GameHandler,
	deviceHandler *DeviceHandler,
	authHandler *AuthHandler,
	wsHub *ws.Hub,
	db *sql.DB,
	mqttClient *mqttclient.Client,
) http.Handler {
	mux := http.NewServeMux()

	// Health check with DB + MQTT ping
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "ok"
		if err := db.Ping(); err != nil {
			dbStatus = "error"
		}
		mqttStatus := "ok"
		if mqttClient != nil && !mqttClient.Ping() {
			mqttStatus = "error"
		}
		status := http.StatusOK
		overall := "ok"
		if dbStatus != "ok" || mqttStatus != "ok" {
			status = http.StatusServiceUnavailable
			overall = "degraded"
		}
		writeJSON(w, status, HealthResponse{
			Status: overall,
			Checks: map[string]HealthCheck{
				"database": {Status: dbStatus},
				"mqtt":     {Status: mqttStatus},
			},
		})
	})

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Auth (public)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// WebSocket (public)
	mux.HandleFunc("GET /ws", wsHub.HandleWS)

	// Player session lookup (public — for mobile app)
	mux.HandleFunc("GET /api/player/session/{code}", gameHandler.GetPlayerSession)

	// Protected routes
	auth := AuthMiddleware(authUC)

	mux.Handle("POST /api/auth/logout", auth(http.HandlerFunc(authHandler.Logout)))

	// Devices
	mux.Handle("GET /api/devices", auth(http.HandlerFunc(deviceHandler.ListAll)))
	mux.Handle("GET /api/devices/available", auth(http.HandlerFunc(deviceHandler.ListAvailable)))

	// Games
	mux.Handle("GET /api/games", auth(http.HandlerFunc(gameHandler.List)))
	mux.Handle("POST /api/games", auth(http.HandlerFunc(gameHandler.Create)))
	mux.Handle("GET /api/games/{id}", auth(http.HandlerFunc(gameHandler.Get)))
	mux.Handle("GET /api/games/{id}/full", auth(http.HandlerFunc(gameHandler.GetFull)))
	mux.Handle("PATCH /api/games/{id}/settings", auth(http.HandlerFunc(gameHandler.UpdateSettings)))
	mux.Handle("POST /api/games/{id}/start", auth(http.HandlerFunc(gameHandler.Start)))
	mux.Handle("POST /api/games/{id}/end", auth(http.HandlerFunc(gameHandler.End)))

	// Teams
	mux.Handle("GET /api/games/{id}/teams", auth(http.HandlerFunc(gameHandler.ListTeams)))
	mux.Handle("POST /api/games/{id}/teams", auth(http.HandlerFunc(gameHandler.AddTeam)))
	mux.Handle("DELETE /api/games/{id}/teams/{teamId}", auth(http.HandlerFunc(gameHandler.RemoveTeam)))

	// Players
	mux.Handle("GET /api/games/{id}/players", auth(http.HandlerFunc(gameHandler.ListPlayers)))
	mux.Handle("POST /api/games/{id}/players", auth(http.HandlerFunc(gameHandler.AddPlayer)))
	mux.Handle("DELETE /api/games/{id}/players/{playerId}", auth(http.HandlerFunc(gameHandler.RemovePlayer)))
	mux.Handle("PATCH /api/games/{id}/players/{playerId}/team", auth(http.HandlerFunc(gameHandler.UpdatePlayerTeam)))

	// Leaderboard & Events
	mux.Handle("GET /api/games/{id}/leaderboard", auth(http.HandlerFunc(gameHandler.Leaderboard)))
	mux.Handle("GET /api/games/{id}/events", auth(http.HandlerFunc(gameHandler.Events)))

	// Apply global middleware (order: outermost first)
	var handler http.Handler = mux
	handler = LoggingMiddleware(handler)
	handler = RequestIDMiddleware(handler)
	handler = RateLimitMiddleware(cfg.RateLimitRPS, cfg.RateLimitBurst)(handler)
	handler = CORSMiddleware(cfg)(handler)

	return handler
}
