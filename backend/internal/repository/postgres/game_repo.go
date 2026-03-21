package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type GameRepo struct {
	db *sql.DB
}

func NewGameRepo(db *sql.DB) *GameRepo {
	return &GameRepo{db: db}
}

func (r *GameRepo) Create(ctx context.Context, game *domain.Game) error {
	configBytes, err := json.Marshal(game.Config)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO games (id, name, status, config_json, created_at, started_at, ended_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		game.ID, game.Name, game.Status, configBytes, game.CreatedAt, game.StartedAt, game.EndedAt,
	)
	return err
}

func (r *GameRepo) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	var g domain.Game
	var configBytes []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, status, config_json, created_at, started_at, ended_at FROM games WHERE id = $1`, id,
	).Scan(&g.ID, &g.Name, &g.Status, &configBytes, &g.CreatedAt, &g.StartedAt, &g.EndedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(configBytes) > 0 {
		_ = json.Unmarshal(configBytes, &g.Config)
	}
	return &g, nil
}

func (r *GameRepo) Update(ctx context.Context, game *domain.Game) error {
	configBytes, err := json.Marshal(game.Config)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE games SET name=$1, status=$2, config_json=$3, started_at=$4, ended_at=$5 WHERE id=$6`,
		game.Name, game.Status, configBytes, game.StartedAt, game.EndedAt, game.ID,
	)
	return err
}

func (r *GameRepo) List(ctx context.Context) ([]*domain.Game, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, status, config_json, created_at, started_at, ended_at FROM games ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*domain.Game
	for rows.Next() {
		var g domain.Game
		var configBytes []byte
		if err := rows.Scan(&g.ID, &g.Name, &g.Status, &configBytes, &g.CreatedAt, &g.StartedAt, &g.EndedAt); err != nil {
			return nil, err
		}
		if len(configBytes) > 0 {
			_ = json.Unmarshal(configBytes, &g.Config)
		}
		games = append(games, &g)
	}
	return games, rows.Err()
}
