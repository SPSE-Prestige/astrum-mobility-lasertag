package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
)

type DeviceUseCase struct {
	devices domain.DeviceRepository
}

func NewDeviceUseCase(devices domain.DeviceRepository) *DeviceUseCase {
	return &DeviceUseCase{devices: devices}
}

// Register is called when an ESP32 sends a registration message via MQTT.
func (uc *DeviceUseCase) Register(ctx context.Context, deviceID string) (*domain.Device, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("register device: %w", domain.ErrValidation)
	}

	d := &domain.Device{
		ID:       uuid.New().String(),
		DeviceID: deviceID,
		Status:   domain.DeviceOnline,
		LastSeen: time.Now(),
	}
	if err := uc.devices.Upsert(ctx, d); err != nil {
		return nil, fmt.Errorf("upsert device %s: %w", deviceID, err)
	}

	result, err := uc.devices.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("get device after register %s: %w", deviceID, err)
	}
	return result, nil
}

// Heartbeat updates the last_seen timestamp for a device.
func (uc *DeviceUseCase) Heartbeat(ctx context.Context, deviceID string) error {
	if err := uc.devices.UpdateLastSeen(ctx, deviceID); err != nil {
		return fmt.Errorf("heartbeat device %s: %w", deviceID, err)
	}
	return nil
}

// MarkOffline sets devices that haven't sent a heartbeat within the timeout as offline.
func (uc *DeviceUseCase) MarkOffline(ctx context.Context, timeout time.Duration) ([]string, error) {
	devices, err := uc.devices.ListByStatus(ctx, domain.DeviceOnline)
	if err != nil {
		return nil, fmt.Errorf("list online devices: %w", err)
	}

	cutoff := time.Now().Add(-timeout)
	var offlineIDs []string
	for _, d := range devices {
		if d.LastSeen.Before(cutoff) {
			if err := uc.devices.UpdateStatus(ctx, d.DeviceID, domain.DeviceOffline); err != nil {
				slog.Error("failed to mark device offline", "device_id", d.DeviceID, "error", err)
				continue
			}
			offlineIDs = append(offlineIDs, d.DeviceID)
		}
	}
	return offlineIDs, nil
}

func (uc *DeviceUseCase) ListAll(ctx context.Context) ([]domain.Device, error) {
	return uc.devices.ListAll(ctx)
}

func (uc *DeviceUseCase) ListAvailable(ctx context.Context) ([]domain.Device, error) {
	return uc.devices.ListByStatus(ctx, domain.DeviceOnline)
}
