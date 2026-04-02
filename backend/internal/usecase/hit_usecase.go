package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
)

type HitUseCase struct {
	games   domain.GameRepository
	players domain.PlayerRepository
	events  domain.EventRepository
	txMgr   domain.TxManager
}

func NewHitUseCase(
	games domain.GameRepository,
	players domain.PlayerRepository,
	events domain.EventRepository,
	txMgr domain.TxManager,
) *HitUseCase {
	return &HitUseCase{games: games, players: players, events: events, txMgr: txMgr}
}

// ProcessHit handles a hit event from MQTT. 1 shot = 1 kill.
// Uses transactions and atomic operations to prevent race conditions.
func (uc *HitUseCase) ProcessHit(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*domain.HitResult, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game: %w", err)
	}
	if game.Status != domain.GameRunning {
		return nil, domain.ErrInvalidGameState
	}

	attacker, err := uc.players.GetByGameAndDevice(ctx, gameID, attackerDeviceID)
	if err != nil {
		return nil, fmt.Errorf("get attacker by device %s: %w", attackerDeviceID, err)
	}
	victim, err := uc.players.GetByGameAndDevice(ctx, gameID, victimDeviceID)
	if err != nil {
		return nil, fmt.Errorf("get victim by device %s: %w", victimDeviceID, err)
	}

	if attacker.ID == victim.ID {
		return nil, domain.ErrSelfHit
	}
	if !attacker.IsAlive {
		return nil, domain.ErrPlayerDead
	}
	if !victim.IsAlive {
		return nil, domain.ErrPlayerDead
	}

	// Friendly fire check
	if !game.Settings.FriendlyFire && attacker.TeamID != nil && victim.TeamID != nil && *attacker.TeamID == *victim.TeamID {
		return nil, domain.ErrFriendlyFire
	}

	scorePerKill := game.Settings.ScorePerKill
	killsPerUpgrade := game.Settings.KillsPerUpgrade
	var result *domain.HitResult

	err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// Atomically kill victim (returns false if already dead = duplicate prevention)
		killed, err := uc.players.KillPlayer(txCtx, victim.ID)
		if err != nil {
			return fmt.Errorf("kill victim: %w", err)
		}
		if !killed {
			return domain.ErrPlayerDead
		}

		// Atomically add kill + score + streak to attacker (handles weapon upgrade)
		streakResult, err := uc.players.AddKillScore(txCtx, attacker.ID, scorePerKill, killsPerUpgrade)
		if err != nil {
			return fmt.Errorf("add kill score: %w", err)
		}

		// Determine if a weapon upgrade occurred
		upgraded := killsPerUpgrade > 0 && streakResult.KillStreak > 0 && streakResult.KillStreak%killsPerUpgrade == 0

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
				"score":              scorePerKill,
			},
			Timestamp: time.Now(),
		}
		if err := uc.events.Create(txCtx, event); err != nil {
			return fmt.Errorf("create kill event: %w", err)
		}

		result = &domain.HitResult{
			Kill:           true,
			AttackerID:     attacker.ID,
			VictimID:       victim.ID,
			AttackerScore:  attacker.Score + scorePerKill,
			AttackerKills:  attacker.Kills + 1,
			WeaponUpgraded: upgraded,
			WeaponLevel:    streakResult.WeaponLevel,
			KillStreak:     streakResult.KillStreak,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// Respawn brings a player back to life.
func (uc *HitUseCase) Respawn(ctx context.Context, gameID, deviceID string) error {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game for respawn: %w", err)
	}
	if game.Status != domain.GameRunning {
		return domain.ErrInvalidGameState
	}

	player, err := uc.players.GetByGameAndDevice(ctx, gameID, deviceID)
	if err != nil {
		return fmt.Errorf("get player for respawn: %w", err)
	}

	if err := uc.players.Respawn(ctx, player.ID); err != nil {
		return fmt.Errorf("respawn player: %w", err)
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
	if err := uc.events.Create(ctx, event); err != nil {
		return fmt.Errorf("create respawn event: %w", err)
	}

	return nil
}
