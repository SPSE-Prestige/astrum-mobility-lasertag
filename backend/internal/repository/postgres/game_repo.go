package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type GameRepo struct{ db *sql.DB }

func NewGameRepo(db *sql.DB) *GameRepo { return &GameRepo{db: db} }

func (r *GameRepo) Create(ctx context.Context, g *domain.Game) error {
	settings, _ := json.Marshal(g.Settings)
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`INSERT INTO games (id, code, status, settings, created_at) VALUES ($1,$2,$3,$4,$5)`,
		g.ID, g.Code, g.Status, settings, g.CreatedAt,
	)
	return err
}

func (r *GameRepo) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	return r.scanGame(getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, code, status, settings, created_at, started_at, ended_at FROM games WHERE id = $1`, id,
	))
}

func (r *GameRepo) GetByCode(ctx context.Context, code string) (*domain.Game, error) {
	return r.scanGame(getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, code, status, settings, created_at, started_at, ended_at FROM games WHERE code = $1`, code,
	))
}

func (r *GameRepo) Update(ctx context.Context, g *domain.Game) error {
	settings, _ := json.Marshal(g.Settings)
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE games SET status=$1, settings=$2, started_at=$3, ended_at=$4 WHERE id=$5`,
		g.Status, settings, g.StartedAt, g.EndedAt, g.ID,
	)
	return err
}

func (r *GameRepo) ListAll(ctx context.Context) ([]domain.Game, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, code, status, settings, created_at, started_at, ended_at FROM games ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanGames(rows)
}

func (r *GameRepo) ListByStatus(ctx context.Context, status domain.GameStatus) ([]domain.Game, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, code, status, settings, created_at, started_at, ended_at FROM games WHERE status=$1 ORDER BY created_at DESC`,
		status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanGames(rows)
}

func (r *GameRepo) scanGame(row *sql.Row) (*domain.Game, error) {
	g := &domain.Game{}
	var settingsJSON []byte
	var startedAt, endedAt sql.NullTime
	err := row.Scan(&g.ID, &g.Code, &g.Status, &settingsJSON, &g.CreatedAt, &startedAt, &endedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &g.Settings)
	if startedAt.Valid {
		t := startedAt.Time
		g.StartedAt = &t
	}
	if endedAt.Valid {
		t := endedAt.Time
		g.EndedAt = &t
	}
	return g, nil
}

func (r *GameRepo) scanGames(rows *sql.Rows) ([]domain.Game, error) {
	var games []domain.Game
	for rows.Next() {
		var g domain.Game
		var settingsJSON []byte
		var startedAt, endedAt sql.NullTime
		if err := rows.Scan(&g.ID, &g.Code, &g.Status, &settingsJSON, &g.CreatedAt, &startedAt, &endedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(settingsJSON, &g.Settings)
		if startedAt.Valid {
			t := startedAt.Time
			g.StartedAt = &t
		}
		if endedAt.Valid {
			t := endedAt.Time
			g.EndedAt = &t
		}
		games = append(games, g)
	}
	return games, rows.Err()
}
