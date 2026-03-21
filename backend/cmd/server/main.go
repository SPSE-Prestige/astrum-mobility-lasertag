package main

import (
	"context"
	"log"
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

//	@title			Laser Tag API
//	@version		1.0
//	@description	REST API for the Aimtec 2026 Laser Tag system — manages games, devices, players, teams and real-time events.
//	@host			localhost:8080
//	@BasePath		/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter "Bearer {token}" (without quotes).
func main() {
	cfg := config.Load()

	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}
	defer container.DB.Close()

	// Connect MQTT
	if err := container.MQTTClient.Connect(); err != nil {
		log.Printf("[WARN] MQTT connection failed: %v (will retry on reconnect)", err)
	}
	defer container.MQTTClient.Disconnect()

	// Background: heartbeat timeout checker (mark devices offline after 30s without heartbeat)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			offlineIDs, err := container.DeviceUC.MarkOffline(ctx, 30*time.Second)
			cancel()
			if err != nil {
				log.Printf("[BG] offline check error: %v", err)
			}
			for _, id := range offlineIDs {
				log.Printf("[BG] device %s marked offline", id)
			}
		}
	}()

	// Background: auto-end games that exceeded their duration
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			games, err := container.GameUC.ListGames(ctx)
			if err != nil {
				cancel()
				continue
			}
			for _, g := range games {
				if g.Status != domain.GameRunning {
					continue
				}
				shouldEnd, _ := container.GameUC.ShouldAutoEnd(ctx, g.ID)
				if shouldEnd {
					game, err := container.GameUC.EndGame(ctx, g.ID)
					if err == nil {
						log.Printf("[BG] game %s auto-ended (duration expired)", game.Code)
					}
				}
			}
			cancel()
		}
	}()

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      container.Handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[HTTP] listening on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTTP] server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[SHUTDOWN] shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[SHUTDOWN] server forced to shutdown: %v", err)
	}
	log.Println("[SHUTDOWN] done")
}
