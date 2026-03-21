package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type GameEventRepo struct {
	db *sql.DB
}

func NewGameEventRepo(db *sql.DB) *GameEventRepo {
	return &GameEventRepo{db: db}
}

func (r *GameEventRepo) Create(ctx context.Context, event *domain.GameEvent) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO game_events (id, game_id, type, player_id, target_id, weapon_id, damage, metadata, timestamp)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		event.ID, event.GameID, event.Type, event.PlayerID, event.TargetID,
		event.WeaponID, event.Damage, event.Metadata, event.Timestamp,
	)
	return err
}

func (r *GameEventRepo) GetByGameID(ctx context.Context, gameID string) ([]*domain.GameEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, type, player_id, target_id, weapon_id, damage, metadata, timestamp
		 FROM game_events WHERE game_id = $1 ORDER BY timestamp ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*domain.GameEvent
	for rows.Next() {
		var e domain.GameEvent
		if err := rows.Scan(&e.ID, &e.GameID, &e.Type, &e.PlayerID, &e.TargetID,
			&e.WeaponID, &e.Damage, &e.Metadata, &e.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, rows.Err()
}
