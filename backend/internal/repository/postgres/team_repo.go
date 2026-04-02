package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type TeamRepo struct{ db *sql.DB }

func NewTeamRepo(db *sql.DB) *TeamRepo { return &TeamRepo{db: db} }

func (r *TeamRepo) Create(ctx context.Context, t *domain.Team) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`INSERT INTO teams (id, game_id, name, color) VALUES ($1,$2,$3,$4)`,
		t.ID, t.GameID, t.Name, t.Color,
	)
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	t := &domain.Team{}
	err := getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, game_id, name, color FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.GameID, &t.Name, &t.Color)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return t, err
}

func (r *TeamRepo) ListByGame(ctx context.Context, gameID string) ([]domain.Team, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, game_id, name, color FROM teams WHERE game_id = $1`, gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.GameID, &t.Name, &t.Color); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *TeamRepo) Delete(ctx context.Context, id string) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM teams WHERE id = $1`, id)
	return err
}
