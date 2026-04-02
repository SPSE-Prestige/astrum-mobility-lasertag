package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain/mocks"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(h)
}

func TestLogin_Success(t *testing.T) {
	password := "secret123"
	user := &domain.User{
		ID:           "user-1",
		Username:     "admin",
		PasswordHash: hashPassword(t, password),
		Role:         domain.RoleAdmin,
	}

	users := &mocks.MockUserRepository{
		GetByUsernameFn: func(_ context.Context, username string) (*domain.User, error) {
			if username == "admin" {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	sessions := &mocks.MockSessionRepository{
		CreateFn: func(_ context.Context, s *domain.Session) error {
			if s.UserID != user.ID {
				t.Errorf("expected UserID %q, got %q", user.ID, s.UserID)
			}
			if s.Token == "" {
				t.Error("expected non-empty token")
			}
			return nil
		},
	}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	session, err := uc.Login(context.Background(), "admin", password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.UserID != user.ID {
		t.Errorf("expected UserID %q, got %q", user.ID, session.UserID)
	}
	if session.Token == "" {
		t.Error("expected non-empty token")
	}
	if session.ExpiresAt.Before(time.Now()) {
		t.Error("expected ExpiresAt in the future")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	users := &mocks.MockUserRepository{
		GetByUsernameFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	sessions := &mocks.MockSessionRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	_, err := uc.Login(context.Background(), "unknown", "pass")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := &domain.User{
		ID:           "user-1",
		Username:     "admin",
		PasswordHash: hashPassword(t, "correct-password"),
	}
	users := &mocks.MockUserRepository{
		GetByUsernameFn: func(_ context.Context, _ string) (*domain.User, error) {
			return user, nil
		},
	}
	sessions := &mocks.MockSessionRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	_, err := uc.Login(context.Background(), "admin", "wrong-password")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestValidateToken_Success(t *testing.T) {
	user := &domain.User{ID: "user-1", Username: "admin", Role: domain.RoleAdmin}
	session := &domain.Session{
		ID:        "sess-1",
		UserID:    user.ID,
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	users := &mocks.MockUserRepository{
		GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			if id == user.ID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	sessions := &mocks.MockSessionRepository{
		GetByTokenFn: func(_ context.Context, token string) (*domain.Session, error) {
			if token == "valid-token" {
				return session, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	got, err := uc.ValidateToken(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("expected user ID %q, got %q", user.ID, got.ID)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	session := &domain.Session{
		ID:        "sess-1",
		UserID:    "user-1",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	deleteByTokenCalled := false
	sessions := &mocks.MockSessionRepository{
		GetByTokenFn: func(_ context.Context, _ string) (*domain.Session, error) {
			return session, nil
		},
		DeleteByTokenFn: func(_ context.Context, token string) error {
			deleteByTokenCalled = true
			if token != "expired-token" {
				t.Errorf("expected token %q, got %q", "expired-token", token)
			}
			return nil
		},
	}
	users := &mocks.MockUserRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	_, err := uc.ValidateToken(context.Background(), "expired-token")
	if !errors.Is(err, domain.ErrSessionExpired) {
		t.Errorf("expected ErrSessionExpired, got %v", err)
	}
	if !deleteByTokenCalled {
		t.Error("expected DeleteByToken to be called for expired session")
	}
}

func TestValidateToken_NotFound(t *testing.T) {
	sessions := &mocks.MockSessionRepository{
		GetByTokenFn: func(_ context.Context, _ string) (*domain.Session, error) {
			return nil, domain.ErrNotFound
		},
	}
	users := &mocks.MockUserRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	_, err := uc.ValidateToken(context.Background(), "nonexistent-token")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogout(t *testing.T) {
	called := false
	sessions := &mocks.MockSessionRepository{
		DeleteByTokenFn: func(_ context.Context, token string) error {
			called = true
			if token != "my-token" {
				t.Errorf("expected token %q, got %q", "my-token", token)
			}
			return nil
		},
	}
	users := &mocks.MockUserRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	if err := uc.Logout(context.Background(), "my-token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected DeleteByToken to be called")
	}
}

func TestCleanupExpiredSessions(t *testing.T) {
	called := false
	sessions := &mocks.MockSessionRepository{
		DeleteExpiredFn: func(_ context.Context) error {
			called = true
			return nil
		},
	}
	users := &mocks.MockUserRepository{}

	uc := NewAuthUseCase(users, sessions, 24*time.Hour)
	if err := uc.CleanupExpiredSessions(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected DeleteExpired to be called")
	}
}
