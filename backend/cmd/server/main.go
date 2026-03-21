package main

import (
	"log"
	"net/http"

	_ "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/docs"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/di"
)

// @title           Laser Tag Game API
// @version         1.0
// @description     Real-time configurable laser tag game engine backend.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}" (without quotes)
func main() {
	cfg := config.Load()
	container := di.NewContainer(cfg)
	defer container.DB.Close()

	addr := ":" + cfg.Server.Port
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, container.Router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
