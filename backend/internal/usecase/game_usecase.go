package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
)

type GameUseCase struct {
	games   domain.GameRepository
	teams   domain.TeamRepository
	players domain.PlayerRepository
	devices domain.DeviceRepository
	events  domain.EventRepository
}

func NewGameUseCase(
	games domain.GameRepository,
	teams domain.TeamRepository,
	players domain.PlayerRepository,
	devices domain.DeviceRepository,
	events domain.EventRepository,
) *GameUseCase {
	return &GameUseCase{
		games:   games,
		teams:   teams,
		players: players,
		devices: devices,
		events:  events,
	}
}

// CreateGame creates a new game in lobby state.
func (uc *GameUseCase) CreateGame(ctx context.Context, settings *domain.GameSettings) (*domain.Game, error) {
	if settings == nil {
		s := domain.DefaultGameSettings()
		settings = &s
	}

	game := &domain.Game{
		ID:        uuid.New().String(),
		Code:      generateGameCode(),
		Status:    domain.GameLobby,
		Settings:  *settings,
		CreatedAt: time.Now(),
	}
	if err := uc.games.Create(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func (uc *GameUseCase) GetGame(ctx context.Context, id string) (*domain.Game, error) {
	return uc.games.GetByID(ctx, id)
}

func (uc *GameUseCase) ListGames(ctx context.Context) ([]domain.Game, error) {
	return uc.games.ListAll(ctx)
}

// AddTeam adds a team to a game.
func (uc *GameUseCase) AddTeam(ctx context.Context, gameID, name, color string) (*domain.Team, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}

	team := &domain.Team{
		ID:     uuid.New().String(),
		GameID: gameID,
		Name:   name,
		Color:  color,
	}
	if err := uc.teams.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *GameUseCase) ListTeams(ctx context.Context, gameID string) ([]domain.Team, error) {
	return uc.teams.ListByGame(ctx, gameID)
}

func (uc *GameUseCase) RemoveTeam(ctx context.Context, teamID string) error {
	return uc.teams.Delete(ctx, teamID)
}

// AddPlayer adds a device as a player to a game. The device must be online.
func (uc *GameUseCase) AddPlayer(ctx context.Context, gameID, deviceID, nickname string, teamID *string) (*domain.Player, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}

	// Check max players
	players, err := uc.players.ListByGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if len(players) >= game.Settings.MaxPlayers {
		return nil, domain.ErrGameFull
	}

	// Check device exists and is available
	device, err := uc.devices.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device.Status == domain.DeviceInGame {
		return nil, domain.ErrDeviceInGame
	}

	player := &domain.Player{
		ID:       uuid.New().String(),
		GameID:   gameID,
		TeamID:   teamID,
		DeviceID: deviceID,
		Nickname: nickname,
		IsAlive:  true,
	}
	if err := uc.players.Create(ctx, player); err != nil {
		return nil, err
	}
	return player, nil
}

func (uc *GameUseCase) RemovePlayer(ctx context.Context, playerID string) error {
	return uc.players.Delete(ctx, playerID)
}

func (uc *GameUseCase) ListPlayers(ctx context.Context, gameID string) ([]domain.Player, error) {
	return uc.players.ListByGame(ctx, gameID)
}

// StartGame transitions game from lobby to running.
func (uc *GameUseCase) StartGame(ctx context.Context, gameID string) (*domain.Game, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}

	now := time.Now()
	game.Status = domain.GameRunning
	game.StartedAt = &now

	if err := uc.games.Update(ctx, game); err != nil {
		return nil, err
	}

	// Mark all player devices as in_game
	players, _ := uc.players.ListByGame(ctx, gameID)
	for _, p := range players {
		_ = uc.devices.UpdateStatus(ctx, p.DeviceID, domain.DeviceInGame)
	}

	return game, nil
}

// EndGame transitions game from running to finished.
func (uc *GameUseCase) EndGame(ctx context.Context, gameID string) (*domain.Game, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameRunning {
		return nil, domain.ErrInvalidGameState
	}

	now := time.Now()
	game.Status = domain.GameFinished
	game.EndedAt = &now

	if err := uc.games.Update(ctx, game); err != nil {
		return nil, err
	}

	// Release devices back to online
	players, _ := uc.players.ListByGame(ctx, gameID)
	for _, p := range players {
		_ = uc.devices.UpdateStatus(ctx, p.DeviceID, domain.DeviceOnline)
	}

	return game, nil
}

func (uc *GameUseCase) UpdateSettings(ctx context.Context, gameID string, settings domain.GameSettings) (*domain.Game, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}
	game.Settings = settings
	if err := uc.games.Update(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func generateGameCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// GetGameFull returns a game with its teams, players, and events.
type GameFull struct {
	Game    domain.Game
	Teams   []domain.Team
	Players []domain.Player
	Events  []domain.GameEvent
}

func (uc *GameUseCase) GetGameFull(ctx context.Context, gameID string) (*GameFull, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	teams, err := uc.teams.ListByGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	players, err := uc.players.ListByGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	events, err := uc.events.ListByGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return &GameFull{
		Game:    *game,
		Teams:   teams,
		Players: players,
		Events:  events,
	}, nil
}

func (uc *GameUseCase) GetLeaderboard(ctx context.Context, gameID string) ([]domain.Player, error) {
	return uc.players.ListByGame(ctx, gameID) // already sorted by score DESC
}

func (uc *GameUseCase) ListEvents(ctx context.Context, gameID string) ([]domain.GameEvent, error) {
	return uc.events.ListByGame(ctx, gameID)
}

func (uc *GameUseCase) UpdatePlayerTeam(ctx context.Context, playerID string, teamID *string) error {
	player, err := uc.players.GetByID(ctx, playerID)
	if err != nil {
		return err
	}
	game, err := uc.games.GetByID(ctx, player.GameID)
	if err != nil {
		return err
	}
	if game.Status != domain.GameLobby {
		return domain.ErrInvalidGameState
	}
	player.TeamID = teamID
	return uc.players.Update(ctx, player)
}

// ShouldAutoEnd checks if the game has exceeded its duration.
func (uc *GameUseCase) ShouldAutoEnd(ctx context.Context, gameID string) (bool, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return false, err
	}
	if game.Status != domain.GameRunning || game.StartedAt == nil || game.Settings.GameDuration == 0 {
		return false, nil
	}
	elapsed := time.Since(*game.StartedAt)
	return elapsed >= time.Duration(game.Settings.GameDuration)*time.Second, nil
}

// RemainingTime returns seconds left.
func (uc *GameUseCase) RemainingTime(game *domain.Game) int {
	if game.Status != domain.GameRunning || game.StartedAt == nil || game.Settings.GameDuration == 0 {
		return -1 // unlimited
	}
	elapsed := time.Since(*game.StartedAt)
	remaining := time.Duration(game.Settings.GameDuration)*time.Second - elapsed
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}
