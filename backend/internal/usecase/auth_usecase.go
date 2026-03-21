package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	users    domain.UserRepository
	sessions domain.SessionRepository
}

func NewAuthUseCase(users domain.UserRepository, sessions domain.SessionRepository) *AuthUseCase {
	return &AuthUseCase{users: users, sessions: sessions}
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*domain.Session, error) {
	user, err := uc.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := uc.sessions.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	session, err := uc.sessions.GetByToken(ctx, token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if time.Now().After(session.ExpiresAt) {
		_ = uc.sessions.DeleteByToken(ctx, token)
		return nil, domain.ErrSessionExpired
	}

	return uc.users.GetByID(ctx, session.UserID)
}

func (uc *AuthUseCase) Logout(ctx context.Context, token string) error {
	return uc.sessions.DeleteByToken(ctx, token)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
