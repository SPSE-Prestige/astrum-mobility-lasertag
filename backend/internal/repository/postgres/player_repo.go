package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type PlayerRepo struct{ db *sql.DB }

func NewPlayerRepo(db *sql.DB) *PlayerRepo { return &PlayerRepo{db: db} }

func (r *PlayerRepo) Create(ctx context.Context, p *domain.Player) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO players (id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.ID, p.GameID, nilIfEmpty(p.TeamID), p.DeviceID, p.Nickname,
		p.Score, p.Kills, p.Deaths, p.IsAlive,
	)
	return err
}

func (r *PlayerRepo) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	return r.scanPlayer(r.db.QueryRowContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive FROM players WHERE id = $1`, id,
	))
}

func (r *PlayerRepo) GetByGameAndDevice(ctx context.Context, gameID, deviceID string) (*domain.Player, error) {
	return r.scanPlayer(r.db.QueryRowContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive
		 FROM players WHERE game_id = $1 AND device_id = $2`, gameID, deviceID,
	))
}

func (r *PlayerRepo) ListByGame(ctx context.Context, gameID string) ([]domain.Player, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive FROM players WHERE game_id = $1 ORDER BY score DESC`, gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPlayers(rows)
}

func (r *PlayerRepo) ListByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive FROM players WHERE team_id = $1 ORDER BY score DESC`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPlayers(rows)
}

func (r *PlayerRepo) Update(ctx context.Context, p *domain.Player) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE players SET team_id=$1, nickname=$2, score=$3, kills=$4, deaths=$5, is_alive=$6 WHERE id=$7`,
		nilIfEmpty(p.TeamID), p.Nickname, p.Score, p.Kills, p.Deaths, p.IsAlive, p.ID,
	)
	return err
}

func (r *PlayerRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM players WHERE id = $1`, id)
	return err
}

func (r *PlayerRepo) scanPlayer(row *sql.Row) (*domain.Player, error) {
	p := &domain.Player{}
	var teamID sql.NullString
	err := row.Scan(&p.ID, &p.GameID, &teamID, &p.DeviceID, &p.Nickname,
		&p.Score, &p.Kills, &p.Deaths, &p.IsAlive)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if teamID.Valid {
		p.TeamID = &teamID.String
	}
	return p, nil
}

func (r *PlayerRepo) scanPlayers(rows *sql.Rows) ([]domain.Player, error) {
	var players []domain.Player
	for rows.Next() {
		var p domain.Player
		var teamID sql.NullString
		if err := rows.Scan(&p.ID, &p.GameID, &teamID, &p.DeviceID, &p.Nickname,
			&p.Score, &p.Kills, &p.Deaths, &p.IsAlive); err != nil {
			return nil, err
		}
		if teamID.Valid {
			p.TeamID = &teamID.String
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func nilIfEmpty(s *string) any {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}
