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
	players domain.PlayerRepository
}

func NewDeviceUseCase(devices domain.DeviceRepository, players domain.PlayerRepository) *DeviceUseCase {
	return &DeviceUseCase{devices: devices, players: players}
}

// Register is called when an ESP32 sends a registration message via MQTT.
func (uc *DeviceUseCase) Register(ctx context.Context, deviceID string) (*domain.Device, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("register device: %w", domain.ErrValidation)
	}

	// Check if device already exists and is in-game — don't overwrite status.
	existing, err := uc.devices.GetByDeviceID(ctx, deviceID)
	if err == nil && existing.Status == domain.DeviceInGame {
		// Device reconnected while in-game — just update last_seen, keep in_game status.
		if err := uc.devices.UpdateLastSeen(ctx, deviceID); err != nil {
			return nil, fmt.Errorf("update last_seen for in-game device %s: %w", deviceID, err)
		}
		existing.LastSeen = time.Now()
		return existing, nil
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

// Reconnect checks if a device is part of a running game and returns its state.
// Returns nil if the device is not in any active game.
func (uc *DeviceUseCase) Reconnect(ctx context.Context, deviceID string) (*domain.ReconnectInfo, error) {
	player, game, err := uc.players.FindActivePlayerByDevice(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("find active game for device %s: %w", deviceID, err)
	}
	if player == nil {
		return nil, nil // not in any active game
	}

	// Restore device status to in_game
	if err := uc.devices.UpdateStatus(ctx, deviceID, domain.DeviceInGame); err != nil {
		slog.Error("failed to restore device status to in_game", "device_id", deviceID, "error", err)
	}

	remaining := -1
	if game.Settings.GameDuration > 0 && game.StartedAt != nil {
		elapsed := int(time.Since(*game.StartedAt).Seconds())
		remaining = game.Settings.GameDuration - elapsed
		if remaining < 0 {
			remaining = 0
		}
	}

	slog.Info("device reconnected to active game",
		"device_id", deviceID,
		"game_id", game.ID,
		"player_id", player.ID,
		"weapon_level", player.WeaponLevel,
		"remaining_time", remaining,
	)

	return &domain.ReconnectInfo{
		Player:        *player,
		Game:          *game,
		RemainingTime: remaining,
	}, nil
}
