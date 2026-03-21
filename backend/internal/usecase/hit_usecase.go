package usecase

import (
	"context"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase/gamemodes"
	"github.com/google/uuid"
)

type HitUseCase struct {
	gameRepo   domain.GameRepository
	playerRepo domain.GamePlayerRepository
	weaponRepo domain.WeaponRepository
	eventRepo  domain.GameEventRepository
	cache      domain.GameCache
	eventBus   domain.EventBus
	registry   *gamemodes.Registry
	gameUC     *GameUseCase
}

func NewHitUseCase(
	gameRepo domain.GameRepository,
	playerRepo domain.GamePlayerRepository,
	weaponRepo domain.WeaponRepository,
	eventRepo domain.GameEventRepository,
	cache domain.GameCache,
	eventBus domain.EventBus,
	registry *gamemodes.Registry,
	gameUC *GameUseCase,
) *HitUseCase {
	return &HitUseCase{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
		weaponRepo: weaponRepo,
		eventRepo:  eventRepo,
		cache:      cache,
		eventBus:   eventBus,
		registry:   registry,
		gameUC:     gameUC,
	}
}

func (uc *HitUseCase) ProcessHit(ctx context.Context, gameID string, hit domain.DeviceHitEvent) error {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameStatusRunning {
		return domain.ErrGameNotRunning
	}

	// Map devices to players
	attacker, err := uc.playerRepo.GetByGunID(ctx, gameID, hit.DeviceID)
	if err != nil {
		return domain.ErrPlayerNotInGame
	}
	victim, err := uc.playerRepo.GetByDeviceID(ctx, gameID, hit.TargetID)
	if err != nil {
		return domain.ErrPlayerNotInGame
	}

	if attacker.ID == victim.ID {
		return domain.ErrSamePlayer
	}

	// Get live states
	attackerState, err := uc.cache.GetPlayerState(ctx, gameID, attacker.ID)
	if err != nil {
		return err
	}
	victimState, err := uc.cache.GetPlayerState(ctx, gameID, victim.ID)
	if err != nil {
		return err
	}

	if !attackerState.IsAlive {
		return domain.ErrPlayerDead
	}
	if !victimState.IsAlive {
		return nil // target already dead, ignore
	}

	// Get game mode handler
	handler, err := uc.registry.Get(game.Config.GameMode)
	if err != nil {
		return err
	}

	// Determine damage
	damage := hit.Damage
	if hit.WeaponID != "" {
		weapon, err := uc.weaponRepo.GetByID(ctx, hit.WeaponID)
		if err == nil {
			damage = weapon.Damage
		}
	}

	// Apply game mode rules
	result := handler.OnHit(game, attacker, victim, damage, hit.IsHeadshot)
	if result.DamageApplied == 0 {
		return nil // blocked (e.g., friendly fire disabled)
	}

	// Update victim state
	victimState.HP -= result.DamageApplied
	if victimState.HP < 0 {
		victimState.HP = 0
	}

	// Update attacker state
	attackerState.Score += result.AttackerScoreChange
	attackerState.Kills += 0 // will be updated on kill

	// Record hit event
	hitEvent := &domain.GameEvent{
		ID:        uuid.New().String(),
		GameID:    gameID,
		Type:      domain.EventTypeHit,
		PlayerID:  attacker.ID,
		TargetID:  victim.ID,
		WeaponID:  hit.WeaponID,
		Damage:    result.DamageApplied,
		Timestamp: time.Now(),
	}
	_ = uc.eventRepo.Create(ctx, hitEvent)

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypeHit),
		GameID: gameID,
		Payload: map[string]interface{}{
			"attacker_id": attacker.ID,
			"victim_id":   victim.ID,
			"damage":      result.DamageApplied,
			"victim_hp":   victimState.HP,
			"is_headshot": hit.IsHeadshot,
		},
	})

	if result.IsKill {
		attackerState.Kills++
		victimState.Deaths++
		victimState.IsAlive = false

		killResult := handler.OnKill(game, attacker, victim)
		attackerState.Score += killResult.AttackerScoreChange

		// Update team score
		if attacker.TeamID != nil {
			_ = uc.cache.IncrTeamScore(ctx, gameID, *attacker.TeamID, game.Config.Scoring.PointsPerKill)
		}

		killEvent := &domain.GameEvent{
			ID:        uuid.New().String(),
			GameID:    gameID,
			Type:      domain.EventTypeKill,
			PlayerID:  attacker.ID,
			TargetID:  victim.ID,
			Timestamp: time.Now(),
		}
		_ = uc.eventRepo.Create(ctx, killEvent)

		uc.eventBus.Publish(gameID, domain.WSMessage{
			Type:   string(domain.EventTypeKill),
			GameID: gameID,
			Payload: map[string]interface{}{
				"attacker_id": attacker.ID,
				"victim_id":   victim.ID,
			},
		})

		// Handle respawn
		if handler.CanRespawn(game, victim) {
			if victimState.LivesRemaining > 0 {
				victimState.LivesRemaining--
			}
			go uc.scheduleRespawn(gameID, victim.ID, game.Config.Player)
		}

		// Check win condition
		players, _ := uc.playerRepo.GetByGameID(ctx, gameID)
		teams, _ := uc.gameUC.GetTeams(ctx, gameID)
		// Update player alive states from cache for win check
		allStates, _ := uc.cache.GetAllPlayerStates(ctx, gameID)
		stateMap := make(map[string]*domain.PlayerLiveState)
		for _, s := range allStates {
			stateMap[s.PlayerID] = s
		}
		for _, p := range players {
			if s, ok := stateMap[p.ID]; ok {
				p.IsAlive = s.IsAlive
				p.LivesRemaining = s.LivesRemaining
			}
		}
		winResult := handler.CheckWinCondition(game, players, teams)
		if winResult.GameOver {
			_ = uc.gameUC.EndGame(ctx, gameID)
		}
	}

	// Persist live states to cache
	_ = uc.cache.SetPlayerState(ctx, gameID, attackerState)
	_ = uc.cache.SetPlayerState(ctx, gameID, victimState)

	return nil
}

func (uc *HitUseCase) scheduleRespawn(gameID, playerID string, playerCfg domain.PlayerConfig) {
	time.Sleep(time.Duration(playerCfg.RespawnDelaySeconds) * time.Second)
	ctx := context.Background()
	state, err := uc.cache.GetPlayerState(ctx, gameID, playerID)
	if err != nil {
		return
	}
	state.IsAlive = true
	state.HP = playerCfg.MaxHP
	_ = uc.cache.SetPlayerState(ctx, gameID, state)

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypeRespawn),
		GameID: gameID,
		Payload: map[string]interface{}{
			"player_id": playerID,
			"hp":        playerCfg.MaxHP,
		},
	})
}
