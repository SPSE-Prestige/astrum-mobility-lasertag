package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	_ "github.com/lib/pq"
)

func NewDB(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return db, nil
}
