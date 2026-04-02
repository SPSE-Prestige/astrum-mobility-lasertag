package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type SessionRepo struct{ db *sql.DB }

func NewSessionRepo(db *sql.DB) *SessionRepo { return &SessionRepo{db: db} }

func (r *SessionRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`INSERT INTO admin_sessions (id, user_id, token, expires_at, created_at) VALUES ($1,$2,$3,$4,$5)`,
		s.ID, s.UserID, s.Token, s.ExpiresAt, s.CreatedAt,
	)
	return err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	s := &domain.Session{}
	err := getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM admin_sessions WHERE token = $1`, token,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return s, err
}

func (r *SessionRepo) DeleteByToken(ctx context.Context, token string) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM admin_sessions WHERE token = $1`, token)
	return err
}

func (r *SessionRepo) DeleteExpired(ctx context.Context) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM admin_sessions WHERE expires_at < NOW()`)
	return err
}
