package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type AdminSessionRepo struct {
	db *sql.DB
}

func NewAdminSessionRepo(db *sql.DB) *AdminSessionRepo {
	return &AdminSessionRepo{db: db}
}

func (r *AdminSessionRepo) Create(ctx context.Context, session *domain.AdminSession) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO admin_sessions (id, user_id, token, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		session.ID, session.UserID, session.Token, session.ExpiresAt, session.CreatedAt,
	)
	return err
}

func (r *AdminSessionRepo) GetByToken(ctx context.Context, token string) (*domain.AdminSession, error) {
	var s domain.AdminSession
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM admin_sessions WHERE token = $1`, token,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &s, err
}

func (r *AdminSessionRepo) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM admin_sessions WHERE user_id = $1`, userID)
	return err
}
