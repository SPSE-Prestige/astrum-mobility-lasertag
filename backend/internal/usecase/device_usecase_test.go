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
			// First call: device not found. Second call: return registered device.
			return expected, nil
		},
	}
	players := &mocks.MockPlayerRepository{}

	uc := NewDeviceUseCase(devices, players)
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
	players := &mocks.MockPlayerRepository{}
	uc := NewDeviceUseCase(devices, players)

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
	players := &mocks.MockPlayerRepository{}

	uc := NewDeviceUseCase(devices, players)
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
	players := &mocks.MockPlayerRepository{}

	uc := NewDeviceUseCase(devices, players)
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
	players := &mocks.MockPlayerRepository{}

	uc := NewDeviceUseCase(devices, players)
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

func TestReconnect_DeviceInGame(t *testing.T) {
	startedAt := time.Now().Add(-60 * time.Second)
	player := &domain.Player{
		ID: "p1", GameID: "g1", DeviceID: "esp32-001",
		Kills: 5, Deaths: 2, Score: 500, IsAlive: true, WeaponLevel: 2, KillStreak: 1,
	}
	game := &domain.Game{
		ID: "g1", Status: domain.GameRunning, StartedAt: &startedAt,
		Settings: domain.GameSettings{GameDuration: 300},
	}

	devices := &mocks.MockDeviceRepository{
		UpdateStatusFn: func(_ context.Context, deviceID string, status domain.DeviceStatus) error {
			if status != domain.DeviceInGame {
				t.Errorf("expected status in_game, got %q", status)
			}
			return nil
		},
	}
	players := &mocks.MockPlayerRepository{
		FindActivePlayerByDeviceFn: func(_ context.Context, deviceID string) (*domain.Player, *domain.Game, error) {
			return player, game, nil
		},
	}

	uc := NewDeviceUseCase(devices, players)
	info, err := uc.Reconnect(context.Background(), "esp32-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected reconnect info, got nil")
	}
	if info.Player.ID != "p1" {
		t.Errorf("expected player ID p1, got %q", info.Player.ID)
	}
	if info.RemainingTime <= 0 || info.RemainingTime > 300 {
		t.Errorf("expected remaining time between 1-300, got %d", info.RemainingTime)
	}
}

func TestReconnect_NotInGame(t *testing.T) {
	devices := &mocks.MockDeviceRepository{}
	players := &mocks.MockPlayerRepository{
		FindActivePlayerByDeviceFn: func(_ context.Context, deviceID string) (*domain.Player, *domain.Game, error) {
			return nil, nil, nil
		},
	}

	uc := NewDeviceUseCase(devices, players)
	info, err := uc.Reconnect(context.Background(), "esp32-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info != nil {
		t.Error("expected nil info for device not in game")
	}
}
