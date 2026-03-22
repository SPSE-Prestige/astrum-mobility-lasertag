package usecase

import (
	"context"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
)

type HitUseCase struct {
	games   domain.GameRepository
	players domain.PlayerRepository
	events  domain.EventRepository
}

func NewHitUseCase(
	games domain.GameRepository,
	players domain.PlayerRepository,
	events domain.EventRepository,
) *HitUseCase {
	return &HitUseCase{games: games, players: players, events: events}
}

// HitResult carries the outcome of a hit event.
type HitResult struct {
	Kill          bool
	AttackerID    string
	VictimID      string
	AttackerScore int
	AttackerKills int
}

// ProcessHit handles a hit event from MQTT. 1 shot = 1 kill.
// Returns ErrPlayerDead if victim/attacker is dead, ErrFriendlyFire if same team.
func (uc *HitUseCase) ProcessHit(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*HitResult, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameRunning {
		return nil, domain.ErrInvalidGameState
	}

	attacker, err := uc.players.GetByGameAndDevice(ctx, gameID, attackerDeviceID)
	if err != nil {
		return nil, err
	}
	victim, err := uc.players.GetByGameAndDevice(ctx, gameID, victimDeviceID)
	if err != nil {
		return nil, err
	}

	// Self-hit check
	if attacker.ID == victim.ID {
		return nil, domain.ErrSelfHit
	}

	// Attacker must be alive
	if !attacker.IsAlive {
		return nil, domain.ErrPlayerDead
	}

	// Victim must be alive (duplicate kill prevention)
	if !victim.IsAlive {
		return nil, domain.ErrPlayerDead
	}

	// Friendly fire check
	if !game.Settings.FriendlyFire && attacker.TeamID != nil && victim.TeamID != nil && *attacker.TeamID == *victim.TeamID {
		return nil, domain.ErrFriendlyFire
	}

	// 1 shot = 1 kill
	victim.IsAlive = false
	victim.Deaths++
	attacker.Kills++
	attacker.Score += 100

	if err := uc.players.Update(ctx, victim); err != nil {
		return nil, err
	}
	if err := uc.players.Update(ctx, attacker); err != nil {
		return nil, err
	}

	// Record event
	event := &domain.GameEvent{
		ID:     uuid.New().String(),
		GameID: gameID,
		Type:   "kill",
		Payload: map[string]any{
			"attacker_id":        attacker.ID,
			"attacker_device_id": attackerDeviceID,
			"attacker_nickname":  attacker.Nickname,
			"victim_id":          victim.ID,
			"victim_device_id":   victimDeviceID,
			"victim_nickname":    victim.Nickname,
		},
		Timestamp: time.Now(),
	}
	_ = uc.events.Create(ctx, event)

	return &HitResult{
		Kill:          true,
		AttackerID:    attacker.ID,
		VictimID:      victim.ID,
		AttackerScore: attacker.Score,
		AttackerKills: attacker.Kills,
	}, nil
}

// Respawn brings a player back to life after the respawn delay.
func (uc *HitUseCase) Respawn(ctx context.Context, gameID, deviceID string) error {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameRunning {
		return domain.ErrInvalidGameState
	}

	player, err := uc.players.GetByGameAndDevice(ctx, gameID, deviceID)
	if err != nil {
		return err
	}

	player.IsAlive = true
	if err := uc.players.Update(ctx, player); err != nil {
		return err
	}

	event := &domain.GameEvent{
		ID:     uuid.New().String(),
		GameID: gameID,
		Type:   "respawn",
		Payload: map[string]any{
			"player_id": player.ID,
			"device_id": deviceID,
			"nickname":  player.Nickname,
		},
		Timestamp: time.Now(),
	}
	_ = uc.events.Create(ctx, event)

	return nil
}
