package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain/mocks"
)

func TestRegister_NewDevice(t *testing.T) {
	expected := &domain.Device{
		ID:       "id-1",
		DeviceID: "esp32-001",
		Status:   domain.DeviceOnline,
		LastSeen: time.Now(),
	}

	devices := &mocks.MockDeviceRepository{
		UpsertFn: func(_ context.Context, d *domain.Device) error {
			if d.DeviceID != "esp32-001" {
				t.Errorf("expected DeviceID 'esp32-001', got %q", d.DeviceID)
			}
			if d.Status != domain.DeviceOnline {
				t.Errorf("expected status Online, got %q", d.Status)
			}
			return nil
		},
		GetByDeviceIDFn: func(_ context.Context, deviceID string) (*domain.Device, error) {
			if deviceID == "esp32-001" {
				return expected, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	uc := NewDeviceUseCase(devices)
	got, err := uc.Register(context.Background(), "esp32-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DeviceID != expected.DeviceID {
		t.Errorf("expected DeviceID %q, got %q", expected.DeviceID, got.DeviceID)
	}
}

func TestRegister_EmptyDeviceID(t *testing.T) {
	devices := &mocks.MockDeviceRepository{}
	uc := NewDeviceUseCase(devices)

	_, err := uc.Register(context.Background(), "")
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestHeartbeat_Success(t *testing.T) {
	called := false
	devices := &mocks.MockDeviceRepository{
		UpdateLastSeenFn: func(_ context.Context, deviceID string) error {
			called = true
			if deviceID != "esp32-001" {
				t.Errorf("expected deviceID 'esp32-001', got %q", deviceID)
			}
			return nil
		},
	}

	uc := NewDeviceUseCase(devices)
	if err := uc.Heartbeat(context.Background(), "esp32-001"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected UpdateLastSeen to be called")
	}
}

func TestListAll(t *testing.T) {
	expected := []domain.Device{
		{ID: "1", DeviceID: "esp-1", Status: domain.DeviceOnline},
		{ID: "2", DeviceID: "esp-2", Status: domain.DeviceOffline},
	}

	devices := &mocks.MockDeviceRepository{
		ListAllFn: func(_ context.Context) ([]domain.Device, error) {
			return expected, nil
		},
	}

	uc := NewDeviceUseCase(devices)
	got, err := uc.ListAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d devices, got %d", len(expected), len(got))
	}
	for i, d := range got {
		if d.DeviceID != expected[i].DeviceID {
			t.Errorf("device[%d]: expected DeviceID %q, got %q", i, expected[i].DeviceID, d.DeviceID)
		}
	}
}

func TestListAvailable(t *testing.T) {
	onlineDevices := []domain.Device{
		{ID: "1", DeviceID: "esp-1", Status: domain.DeviceOnline},
		{ID: "3", DeviceID: "esp-3", Status: domain.DeviceOnline},
	}

	devices := &mocks.MockDeviceRepository{
		ListByStatusFn: func(_ context.Context, status domain.DeviceStatus) ([]domain.Device, error) {
			if status != domain.DeviceOnline {
				t.Errorf("expected status %q, got %q", domain.DeviceOnline, status)
			}
			return onlineDevices, nil
		},
	}

	uc := NewDeviceUseCase(devices)
	got, err := uc.ListAvailable(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 available devices, got %d", len(got))
	}
	for _, d := range got {
		if d.Status != domain.DeviceOnline {
			t.Errorf("expected all devices to be online, got %q", d.Status)
		}
	}
}
