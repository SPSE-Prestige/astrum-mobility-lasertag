package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/di"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"

	_ "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/docs"
)

// @title			Laser Tag API
// @version		1.0
// @description	REST API for the Aimtec 2026 Laser Tag system — manages games, devices, players, teams and real-time events.
// @host			localhost:8080
// @BasePath		/api
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Enter "Bearer {token}" (without quotes).
func main() {
	// Structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	container, err := di.NewContainer(cfg)
	if err != nil {
		slog.Error("failed to initialize", "error", err)
		os.Exit(1)
	}
	defer container.DB.Close()

	// Connect MQTT
	if err := container.MQTTClient.Connect(); err != nil {
		slog.Warn("MQTT connection failed, will retry on reconnect", "error", err)
	}
	defer container.MQTTClient.Disconnect()

	// Context for background goroutines
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Background: heartbeat timeout checker
	go runBackgroundTask(ctx, cfg.HeartbeatCheckInterval, func() {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		offlineIDs, err := container.DeviceUC.MarkOffline(tctx, cfg.DeviceOfflineTimeout)
		if err != nil {
			slog.Error("offline check error", "error", err)
		}
		for _, id := range offlineIDs {
			slog.Info("device marked offline", "device_id", id)
		}
	})

	// Background: auto-end games that exceeded their duration
	go runBackgroundTask(ctx, cfg.AutoEndCheckInterval, func() {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		games, err := container.GameUC.ListGames(tctx)
		if err != nil {
			return
		}
		for _, g := range games {
			if g.Status != domain.GameRunning {
				continue
			}
			shouldEnd, _ := container.GameUC.ShouldAutoEnd(tctx, g.ID)
			if shouldEnd {
				_, _, endErr := container.GameUC.EndGame(tctx, g.ID)
				if endErr == nil {
					slog.Info("game auto-ended", "game_code", g.Code, "game_id", g.ID)
				}
			}
		}
	})

	// Background: session cleanup
	go runBackgroundTask(ctx, cfg.SessionCleanupInterval, func() {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := container.AuthUC.CleanupExpiredSessions(tctx); err != nil {
			slog.Error("session cleanup error", "error", err)
		}
	})

	// HTTP server with configurable timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      container.Handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		slog.Info("HTTP server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("shutdown complete")
}

// runBackgroundTask runs fn on an interval until ctx is cancelled.
func runBackgroundTask(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}
