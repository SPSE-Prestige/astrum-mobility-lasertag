package di

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	httpdelivery "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/http"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/delivery/ws"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/infrastructure/eventbus"
	rediscache "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/infrastructure/redis"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/repository/postgres"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase/gamemodes"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	DB          *sql.DB
	RedisClient *redis.Client
	Router      http.Handler
}

func NewContainer(cfg *config.Config) *Container {
	// Database
	db, err := postgres.NewDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Repositories
	userRepo := postgres.NewUserRepo(db)
	gameRepo := postgres.NewGameRepo(db)
	playerRepo := postgres.NewGamePlayerRepo(db)
	teamRepo := postgres.NewTeamRepo(db)
	weaponRepo := postgres.NewWeaponRepo(db)
	eventRepo := postgres.NewGameEventRepo(db)
	sessionRepo := postgres.NewAdminSessionRepo(db)

	// Infrastructure
	cache := rediscache.NewGameCache(redisClient)
	bus := eventbus.New()

	// Game modes
	registry := gamemodes.NewRegistry()

	// Usecases
	gameUC := usecase.NewGameUseCase(gameRepo, playerRepo, teamRepo, eventRepo, cache, bus, registry)
	hitUC := usecase.NewHitUseCase(gameRepo, playerRepo, weaponRepo, eventRepo, cache, bus, registry, gameUC)
	adminUC := usecase.NewAdminUseCase(userRepo, sessionRepo, gameUC, playerRepo, cache, bus, cfg.JWT)
	playerUC := usecase.NewPlayerUseCase(playerRepo, gameRepo, teamRepo, eventRepo, cache, bus)

	// Delivery - HTTP
	adminHandler := httpdelivery.NewAdminHandler(adminUC)
	gameHandler := httpdelivery.NewGameHandler(gameUC, adminUC)
	playerHandler := httpdelivery.NewPlayerHandler(playerUC)
	deviceHandler := httpdelivery.NewDeviceHandler(hitUC, gameRepo)
	authMiddleware := httpdelivery.NewAuthMiddleware(adminUC)

	// Delivery - WebSocket
	hub := ws.NewHub(bus)
	wsHandler := ws.NewHandler(hub)

	router := httpdelivery.NewRouter(adminHandler, gameHandler, playerHandler, deviceHandler, authMiddleware, wsHandler)

	return &Container{
		DB:          db,
		RedisClient: redisClient,
		Router:      router,
	}
}
