package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type EventRepo struct{ db *sql.DB }

func NewEventRepo(db *sql.DB) *EventRepo { return &EventRepo{db: db} }

func (r *EventRepo) Create(ctx context.Context, e *domain.GameEvent) error {
	payload, _ := json.Marshal(e.Payload)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO game_events (id, game_id, type, payload, timestamp) VALUES ($1,$2,$3,$4,$5)`,
		e.ID, e.GameID, e.Type, payload, e.Timestamp,
	)
	return err
}

func (r *EventRepo) ListByGame(ctx context.Context, gameID string) ([]domain.GameEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, type, payload, timestamp FROM game_events WHERE game_id = $1 ORDER BY timestamp ASC`,
		gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.GameEvent
	for rows.Next() {
		var e domain.GameEvent
		var payloadJSON []byte
		if err := rows.Scan(&e.ID, &e.GameID, &e.Type, &payloadJSON, &e.Timestamp); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(payloadJSON, &e.Payload)
		events = append(events, e)
	}
	return events, rows.Err()
}
