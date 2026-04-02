package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain/mocks"
)

func TestLogin_Success(t *testing.T) {
	expires := time.Now().Add(time.Hour).Truncate(time.Second)
	mock := &mocks.MockAuthUseCasePort{
		LoginFn: func(_ context.Context, username, password string) (*domain.Session, error) {
			if username != "admin" || password != "secret" {
				t.Fatalf("unexpected credentials: %s / %s", username, password)
			}
			return &domain.Session{Token: "tok123", ExpiresAt: expires}, nil
		},
	}
	h := NewAuthHandler(mock)

	body, _ := json.Marshal(LoginRequest{Username: "admin", Password: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp LoginResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token != "tok123" {
		t.Errorf("expected token tok123, got %s", resp.Token)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := NewAuthHandler(&mocks.MockAuthUseCasePort{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "BAD_REQUEST" {
		t.Errorf("expected code BAD_REQUEST, got %s", errResp.Error.Code)
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	mock := &mocks.MockAuthUseCasePort{
		LoginFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
			return nil, domain.ErrUnauthorized
		},
	}
	h := NewAuthHandler(mock)

	body, _ := json.Marshal(LoginRequest{Username: "admin", Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "UNAUTHORIZED" {
		t.Errorf("expected code UNAUTHORIZED, got %s", errResp.Error.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	h := NewAuthHandler(&mocks.MockAuthUseCasePort{})

	body, _ := json.Marshal(LoginRequest{Username: "", Password: ""})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", errResp.Error.Code)
	}
}

func TestLogout_Success(t *testing.T) {
	var logoutToken string
	mock := &mocks.MockAuthUseCasePort{
		LogoutFn: func(_ context.Context, token string) error {
			logoutToken = token
			return nil
		},
	}
	h := NewAuthHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer tok123")
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if logoutToken != "tok123" {
		t.Errorf("expected logout token tok123, got %s", logoutToken)
	}
}
