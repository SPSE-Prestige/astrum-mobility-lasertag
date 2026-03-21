package usecase

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase/gamemodes"
	"github.com/google/uuid"
)

type GameUseCase struct {
	gameRepo   domain.GameRepository
	playerRepo domain.GamePlayerRepository
	teamRepo   domain.TeamRepository
	eventRepo  domain.GameEventRepository
	cache      domain.GameCache
	eventBus   domain.EventBus
	registry   *gamemodes.Registry

	mu     sync.Mutex
	timers map[string]context.CancelFunc
	paused map[string]int // remaining seconds when paused
}

func NewGameUseCase(
	gameRepo domain.GameRepository,
	playerRepo domain.GamePlayerRepository,
	teamRepo domain.TeamRepository,
	eventRepo domain.GameEventRepository,
	cache domain.GameCache,
	eventBus domain.EventBus,
	registry *gamemodes.Registry,
) *GameUseCase {
	return &GameUseCase{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
		eventRepo:  eventRepo,
		cache:      cache,
		eventBus:   eventBus,
		registry:   registry,
		timers:     make(map[string]context.CancelFunc),
		paused:     make(map[string]int),
	}
}

func (uc *GameUseCase) CreateGame(ctx context.Context, name string, config domain.GameConfig) (*domain.Game, error) {
	if _, err := uc.registry.Get(config.GameMode); err != nil {
		return nil, err
	}
	game := &domain.Game{
		ID:        uuid.New().String(),
		Name:      name,
		Status:    domain.GameStatusPending,
		Config:    config,
		CreatedAt: time.Now(),
	}
	if err := uc.gameRepo.Create(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func (uc *GameUseCase) GetGame(ctx context.Context, id string) (*domain.Game, error) {
	return uc.gameRepo.GetByID(ctx, id)
}

func (uc *GameUseCase) ListGames(ctx context.Context) ([]*domain.Game, error) {
	return uc.gameRepo.List(ctx)
}

func (uc *GameUseCase) CreateTeam(ctx context.Context, gameID, name, color string) (*domain.Team, error) {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameStatusPending {
		return nil, domain.ErrGameNotPending
	}
	team := &domain.Team{
		ID:     uuid.New().String(),
		GameID: gameID,
		Name:   name,
		Color:  color,
	}
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *GameUseCase) GetTeams(ctx context.Context, gameID string) ([]*domain.Team, error) {
	return uc.teamRepo.GetByGameID(ctx, gameID)
}

func (uc *GameUseCase) StartGame(ctx context.Context, gameID string) error {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameStatusPending {
		return domain.ErrGameNotPending
	}

	now := time.Now()
	game.Status = domain.GameStatusRunning
	game.StartedAt = &now
	if err := uc.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	// Initialize live state
	players, _ := uc.playerRepo.GetByGameID(ctx, gameID)
	teams, _ := uc.teamRepo.GetByGameID(ctx, gameID)
	teamScores := make(map[string]int)
	for _, t := range teams {
		teamScores[t.ID] = 0
	}
	_ = uc.cache.SetGameState(ctx, &domain.GameLiveState{
		GameID:         gameID,
		Status:         domain.GameStatusRunning,
		TimeRemainingS: game.Config.DurationSeconds,
		TeamScores:     teamScores,
	})
	for _, p := range players {
		_ = uc.cache.SetPlayerState(ctx, gameID, &domain.PlayerLiveState{
			PlayerID:       p.ID,
			GameID:         gameID,
			HP:             p.HP,
			Score:          0,
			Kills:          0,
			Deaths:         0,
			IsAlive:        true,
			LivesRemaining: p.LivesRemaining,
		})
	}

	// Record event
	uc.recordEvent(ctx, gameID, domain.EventTypeGameStart, "", "", "", 0, nil)

	// Start game timer
	uc.startTimer(gameID, game.Config.DurationSeconds)

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypeGameStart),
		GameID: gameID,
		Payload: map[string]interface{}{
			"game": game,
		},
	})

	return nil
}

func (uc *GameUseCase) PauseGame(ctx context.Context, gameID string) error {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameStatusRunning {
		return domain.ErrGameNotRunning
	}

	game.Status = domain.GameStatusPaused
	if err := uc.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	// Cancel timer, store remaining time
	uc.mu.Lock()
	if cancel, ok := uc.timers[gameID]; ok {
		cancel()
		delete(uc.timers, gameID)
	}
	uc.mu.Unlock()

	state, _ := uc.cache.GetGameState(ctx, gameID)
	if state != nil {
		uc.mu.Lock()
		uc.paused[gameID] = state.TimeRemainingS
		uc.mu.Unlock()
		state.Status = domain.GameStatusPaused
		_ = uc.cache.SetGameState(ctx, state)
	}

	uc.recordEvent(ctx, gameID, domain.EventTypeGamePause, "", "", "", 0, nil)
	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type: string(domain.EventTypeGamePause), GameID: gameID,
		Payload: map[string]interface{}{"status": "paused"},
	})
	return nil
}

func (uc *GameUseCase) ResumeGame(ctx context.Context, gameID string) error {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameStatusPaused {
		return domain.ErrGameNotPaused
	}

	game.Status = domain.GameStatusRunning
	if err := uc.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	uc.mu.Lock()
	remaining, ok := uc.paused[gameID]
	if !ok {
		remaining = game.Config.DurationSeconds
	}
	delete(uc.paused, gameID)
	uc.mu.Unlock()

	state, _ := uc.cache.GetGameState(ctx, gameID)
	if state != nil {
		state.Status = domain.GameStatusRunning
		_ = uc.cache.SetGameState(ctx, state)
	}

	uc.startTimer(gameID, remaining)

	uc.recordEvent(ctx, gameID, domain.EventTypeGameResume, "", "", "", 0, nil)
	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type: string(domain.EventTypeGameResume), GameID: gameID,
		Payload: map[string]interface{}{"status": "running", "time_remaining_s": remaining},
	})
	return nil
}

func (uc *GameUseCase) EndGame(ctx context.Context, gameID string) error {
	game, err := uc.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game.Status == domain.GameStatusFinished {
		return nil
	}

	now := time.Now()
	game.Status = domain.GameStatusFinished
	game.EndedAt = &now
	if err := uc.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	// Cancel timer
	uc.mu.Lock()
	if cancel, ok := uc.timers[gameID]; ok {
		cancel()
		delete(uc.timers, gameID)
	}
	delete(uc.paused, gameID)
	uc.mu.Unlock()

	// Persist final player states from cache to DB
	states, _ := uc.cache.GetAllPlayerStates(ctx, gameID)
	for _, s := range states {
		p, err := uc.playerRepo.GetByID(ctx, s.PlayerID)
		if err != nil {
			continue
		}
		p.HP = s.HP
		p.Score = s.Score
		p.Kills = s.Kills
		p.Deaths = s.Deaths
		p.IsAlive = s.IsAlive
		p.LivesRemaining = s.LivesRemaining
		_ = uc.playerRepo.Update(ctx, p)
	}

	_ = uc.cache.DeleteGameState(ctx, gameID)

	uc.recordEvent(ctx, gameID, domain.EventTypeGameEnd, "", "", "", 0, nil)
	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type: string(domain.EventTypeGameEnd), GameID: gameID,
		Payload: map[string]interface{}{"game": game},
	})
	return nil
}

func (uc *GameUseCase) GetGameState(ctx context.Context, gameID string) (*domain.GameLiveState, error) {
	return uc.cache.GetGameState(ctx, gameID)
}

func (uc *GameUseCase) startTimer(gameID string, durationSeconds int) {
	ctx, cancel := context.WithCancel(context.Background())
	uc.mu.Lock()
	uc.timers[gameID] = cancel
	uc.mu.Unlock()

	go func() {
		remaining := durationSeconds
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				remaining--
				_ = uc.cache.SetTimeRemaining(context.Background(), gameID, remaining)
				if remaining <= 0 {
					_ = uc.EndGame(context.Background(), gameID)
					return
				}
			}
		}
	}()
}

func (uc *GameUseCase) recordEvent(ctx context.Context, gameID string, eventType domain.EventType, playerID, targetID, weaponID string, damage int, metadata map[string]interface{}) {
	var metaJSON json.RawMessage
	if metadata != nil {
		metaJSON, _ = json.Marshal(metadata)
	}
	event := &domain.GameEvent{
		ID:        uuid.New().String(),
		GameID:    gameID,
		Type:      eventType,
		PlayerID:  playerID,
		TargetID:  targetID,
		WeaponID:  weaponID,
		Damage:    damage,
		Metadata:  metaJSON,
		Timestamp: time.Now(),
	}
	_ = uc.eventRepo.Create(ctx, event)
}
