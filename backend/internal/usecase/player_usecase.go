package usecase

import (
	"context"
	"sort"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
)

type PlayerUseCase struct {
	playerRepo domain.GamePlayerRepository
	gameRepo   domain.GameRepository
	teamRepo   domain.TeamRepository
	eventRepo  domain.GameEventRepository
	cache      domain.GameCache
	eventBus   domain.EventBus
}

func NewPlayerUseCase(
	playerRepo domain.GamePlayerRepository,
	gameRepo domain.GameRepository,
	teamRepo domain.TeamRepository,
	eventRepo domain.GameEventRepository,
	cache domain.GameCache,
	eventBus domain.EventBus,
) *PlayerUseCase {
	return &PlayerUseCase{
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
		teamRepo:   teamRepo,
		eventRepo:  eventRepo,
		cache:      cache,
		eventBus:   eventBus,
	}
}

func (uc *PlayerUseCase) JoinGame(ctx context.Context, gameID string, nickname, deviceID, gunID string, userID *string, teamID *string, weaponID *string) (*domain.GamePlayer, error) {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameStatusPending {
		return nil, domain.ErrGameNotPending
	}

	count, err := uc.playerRepo.CountByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if count >= game.Config.MaxPlayers {
		return nil, domain.ErrGameFull
	}

	player := &domain.GamePlayer{
		ID:             uuid.New().String(),
		GameID:         gameID,
		UserID:         userID,
		TeamID:         teamID,
		Nickname:       nickname,
		DeviceID:       deviceID,
		GunID:          gunID,
		WeaponID:       weaponID,
		HP:             game.Config.Player.MaxHP,
		Score:          0,
		Kills:          0,
		Deaths:         0,
		IsAlive:        true,
		LivesRemaining: game.Config.Player.Lives,
	}

	if err := uc.playerRepo.Create(ctx, player); err != nil {
		return nil, err
	}

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypePlayerJoin),
		GameID: gameID,
		Payload: map[string]interface{}{
			"player": player,
		},
	})

	return player, nil
}

func (uc *PlayerUseCase) GetLeaderboard(ctx context.Context, gameID string) ([]*domain.LeaderboardEntry, error) {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	players, err := uc.playerRepo.GetByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	teams, _ := uc.teamRepo.GetByGameID(ctx, gameID)
	teamMap := make(map[string]string)
	for _, t := range teams {
		teamMap[t.ID] = t.Name
	}

	entries := make([]*domain.LeaderboardEntry, 0, len(players))

	// If game is running, use live state from cache
	if game.Status == domain.GameStatusRunning || game.Status == domain.GameStatusPaused {
		states, _ := uc.cache.GetAllPlayerStates(ctx, gameID)
		stateMap := make(map[string]*domain.PlayerLiveState)
		for _, s := range states {
			stateMap[s.PlayerID] = s
		}

		// Calculate damage dealt from events
		events, _ := uc.eventRepo.GetByGameID(ctx, gameID)
		damageMap := make(map[string]int)
		for _, e := range events {
			if e.Type == domain.EventTypeHit {
				damageMap[e.PlayerID] += e.Damage
			}
		}

		for _, p := range players {
			entry := &domain.LeaderboardEntry{
				PlayerID:    p.ID,
				Nickname:    p.Nickname,
				TeamID:      p.TeamID,
				DamageDealt: damageMap[p.ID],
			}
			if p.TeamID != nil {
				entry.TeamName = teamMap[*p.TeamID]
			}
			if s, ok := stateMap[p.ID]; ok {
				entry.Score = s.Score
				entry.Kills = s.Kills
				entry.Deaths = s.Deaths
			}
			entries = append(entries, entry)
		}
	} else {
		for _, p := range players {
			entry := &domain.LeaderboardEntry{
				PlayerID: p.ID,
				Nickname: p.Nickname,
				TeamID:   p.TeamID,
				Score:    p.Score,
				Kills:    p.Kills,
				Deaths:   p.Deaths,
			}
			if p.TeamID != nil {
				entry.TeamName = teamMap[*p.TeamID]
			}
			entries = append(entries, entry)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	return entries, nil
}

func (uc *PlayerUseCase) GetPlayers(ctx context.Context, gameID string) ([]*domain.GamePlayer, error) {
	return uc.playerRepo.GetByGameID(ctx, gameID)
}
