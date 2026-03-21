package di

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	httpdelivery "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/http"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/ws"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/infrastructure/mqtt"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/repository/postgres"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

// Container holds all wired dependencies.
type Container struct {
	DB         *sql.DB
	WSHub      *ws.Hub
	MQTTClient *mqtt.Client
	Handler    http.Handler
	Config     *config.Config

	// UseCases (exposed for background tasks)
	DeviceUC *usecase.DeviceUseCase
	GameUC   *usecase.GameUseCase
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Database
	db, err := postgres.NewDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)
	if err != nil {
		return nil, err
	}
	log.Println("[DI] database connected")

	// Repositories
	userRepo := postgres.NewUserRepo(db)
	sessionRepo := postgres.NewSessionRepo(db)
	deviceRepo := postgres.NewDeviceRepo(db)
	gameRepo := postgres.NewGameRepo(db)
	teamRepo := postgres.NewTeamRepo(db)
	playerRepo := postgres.NewPlayerRepo(db)
	eventRepo := postgres.NewEventRepo(db)

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, sessionRepo)
	deviceUC := usecase.NewDeviceUseCase(deviceRepo)
	gameUC := usecase.NewGameUseCase(gameRepo, teamRepo, playerRepo, deviceRepo, eventRepo)
	hitUC := usecase.NewHitUseCase(gameRepo, playerRepo, eventRepo)

	// WebSocket hub
	wsHub := ws.NewHub()

	// MQTT client
	mqttClient := mqtt.NewClient(cfg.MQTTBroker, deviceUC, hitUC, gameUC, wsHub.Broadcast)

	// HTTP handlers
	authHandler := httpdelivery.NewAuthHandler(authUC)
	gameHandler := httpdelivery.NewGameHandler(gameUC, mqttClient)
	deviceHandler := httpdelivery.NewDeviceHandler(deviceUC)

	// Router
	handler := httpdelivery.NewRouter(authUC, gameHandler, deviceHandler, authHandler, wsHub)

	return &Container{
		DB:         db,
		WSHub:      wsHub,
		MQTTClient: mqttClient,
		Handler:    handler,
		Config:     cfg,
		DeviceUC:   deviceUC,
		GameUC:     gameUC,
	}, nil
}
