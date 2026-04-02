package di

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	httpdelivery "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/http"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/ws"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
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

	// Use case ports for background tasks
	AuthUC   domain.AuthUseCasePort
	DeviceUC domain.DeviceUseCasePort
	GameUC   domain.GameUseCasePort
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Database
	db, err := postgres.NewDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)
	if err != nil {
		return nil, err
	}
	slog.Info("database connected")

	txMgr := postgres.NewTxManager(db)

	// Repositories
	userRepo := postgres.NewUserRepo(db)
	sessionRepo := postgres.NewSessionRepo(db)
	deviceRepo := postgres.NewDeviceRepo(db)
	gameRepo := postgres.NewGameRepo(db)
	teamRepo := postgres.NewTeamRepo(db)
	playerRepo := postgres.NewPlayerRepo(db)
	eventRepo := postgres.NewEventRepo(db)

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, sessionRepo, cfg.SessionTTL)
	deviceUC := usecase.NewDeviceUseCase(deviceRepo, playerRepo)
	gameUC := usecase.NewGameUseCase(gameRepo, teamRepo, playerRepo, deviceRepo, eventRepo, txMgr)
	hitUC := usecase.NewHitUseCase(gameRepo, playerRepo, eventRepo, txMgr)

	// WebSocket hub
	wsHub := ws.NewHub()

	// MQTT client (depends on port interfaces, not concrete types)
	mqttClient := mqtt.NewClient(cfg.MQTTBroker, deviceUC, hitUC, gameUC, wsHub)

	// HTTP handlers (depend on port interfaces)
	authHandler := httpdelivery.NewAuthHandler(authUC)
	gameHandler := httpdelivery.NewGameHandler(gameUC, mqttClient)
	deviceHandler := httpdelivery.NewDeviceHandler(deviceUC)

	// Router
	handler := httpdelivery.NewRouter(cfg, authUC, gameHandler, deviceHandler, authHandler, wsHub, db, mqttClient)

	return &Container{
		DB:         db,
		WSHub:      wsHub,
		MQTTClient: mqttClient,
		Handler:    handler,
		Config:     cfg,
		AuthUC:     authUC,
		DeviceUC:   deviceUC,
		GameUC:     gameUC,
	}, nil
}
