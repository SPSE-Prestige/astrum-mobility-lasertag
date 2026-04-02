package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain/mocks"
)

func TestListAll_Success(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	devices := []domain.Device{
		{ID: "1", DeviceID: "dev-a", Status: domain.DeviceOnline, LastSeen: now},
		{ID: "2", DeviceID: "dev-b", Status: domain.DeviceOffline, LastSeen: now},
	}
	mock := &mocks.MockDeviceUseCasePort{
		ListAllFn: func(_ context.Context) ([]domain.Device, error) {
			return devices, nil
		},
	}
	h := NewDeviceHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rec := httptest.NewRecorder()

	h.ListAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp []DeviceResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(resp))
	}
	if resp[0].DeviceID != "dev-a" {
		t.Errorf("expected dev-a, got %s", resp[0].DeviceID)
	}
	if resp[0].Status != "online" {
		t.Errorf("expected online, got %s", resp[0].Status)
	}
	if resp[1].Status != "offline" {
		t.Errorf("expected offline, got %s", resp[1].Status)
	}
}

func TestListAll_Error(t *testing.T) {
	mock := &mocks.MockDeviceUseCasePort{
		ListAllFn: func(_ context.Context) ([]domain.Device, error) {
			return nil, domain.ErrInternal
		},
	}
	h := NewDeviceHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rec := httptest.NewRecorder()

	h.ListAll(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestListAvailable_Success(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	devices := []domain.Device{
		{ID: "1", DeviceID: "dev-a", Status: domain.DeviceOnline, LastSeen: now},
	}
	mock := &mocks.MockDeviceUseCasePort{
		ListAvailableFn: func(_ context.Context) ([]domain.Device, error) {
			return devices, nil
		},
	}
	h := NewDeviceHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/devices/available", nil)
	rec := httptest.NewRecorder()

	h.ListAvailable(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp []DeviceResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 device, got %d", len(resp))
	}
	if resp[0].DeviceID != "dev-a" {
		t.Errorf("expected dev-a, got %s", resp[0].DeviceID)
	}
}
