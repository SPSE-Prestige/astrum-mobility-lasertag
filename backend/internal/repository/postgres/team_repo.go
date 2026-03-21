package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team *domain.Team) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO teams (id, game_id, name, color) VALUES ($1, $2, $3, $4)`,
		team.ID, team.GameID, team.Name, team.Color,
	)
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	var t domain.Team
	err := r.db.QueryRowContext(ctx,
		`SELECT id, game_id, name, color FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.GameID, &t.Name, &t.Color)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &t, err
}

func (r *TeamRepo) GetByGameID(ctx context.Context, gameID string) ([]*domain.Team, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, name, color FROM teams WHERE game_id = $1`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []*domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.GameID, &t.Name, &t.Color); err != nil {
			return nil, err
		}
		teams = append(teams, &t)
	}
	return teams, rows.Err()
}
