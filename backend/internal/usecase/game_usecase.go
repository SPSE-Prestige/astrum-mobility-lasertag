package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	txMgr   domain.TxManager
}

func NewGameUseCase(
	games domain.GameRepository,
	teams domain.TeamRepository,
	players domain.PlayerRepository,
	devices domain.DeviceRepository,
	events domain.EventRepository,
	txMgr domain.TxManager,
) *GameUseCase {
	return &GameUseCase{
		games:   games,
		teams:   teams,
		players: players,
		devices: devices,
		events:  events,
		txMgr:   txMgr,
	}
}

const maxGameCodeRetries = 5

// CreateGame creates a new game in lobby state with unique game code.
func (uc *GameUseCase) CreateGame(ctx context.Context, settings *domain.GameSettings) (*domain.Game, error) {
	if settings == nil {
		s := domain.DefaultGameSettings()
		settings = &s
	}

	if err := settings.Validate(); err != nil {
		return nil, fmt.Errorf("validate settings: %w", err)
	}

	// Retry loop for unique game code
	for attempt := 0; attempt < maxGameCodeRetries; attempt++ {
		game := &domain.Game{
			ID:        uuid.New().String(),
			Code:      generateGameCode(),
			Status:    domain.GameLobby,
			Settings:  *settings,
			CreatedAt: time.Now(),
		}
		err := uc.games.Create(ctx, game)
		if err == nil {
			return game, nil
		}
		// If it's a unique constraint violation, retry with new code
		if attempt < maxGameCodeRetries-1 {
			slog.Warn("game code collision, retrying", "code", game.Code, "attempt", attempt+1)
			continue
		}
		return nil, fmt.Errorf("create game after %d attempts: %w", maxGameCodeRetries, err)
	}
	return nil, fmt.Errorf("create game: exhausted retries")
}

func (uc *GameUseCase) GetGame(ctx context.Context, id string) (*domain.Game, error) {
	game, err := uc.games.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get game %s: %w", id, err)
	}
	return game, nil
}

func (uc *GameUseCase) ListGames(ctx context.Context) ([]domain.Game, error) {
	return uc.games.ListAll(ctx)
}

// AddTeam adds a team to a game (lobby only).
func (uc *GameUseCase) AddTeam(ctx context.Context, gameID, name, color string) (*domain.Team, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game for add team: %w", err)
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
		return nil, fmt.Errorf("create team: %w", err)
	}
	return team, nil
}

func (uc *GameUseCase) ListTeams(ctx context.Context, gameID string) ([]domain.Team, error) {
	return uc.teams.ListByGame(ctx, gameID)
}

// RemoveTeam removes a team only if the game is in lobby state.
func (uc *GameUseCase) RemoveTeam(ctx context.Context, gameID, teamID string) error {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game for remove team: %w", err)
	}
	if game.Status != domain.GameLobby {
		return domain.ErrInvalidGameState
	}
	return uc.teams.Delete(ctx, teamID)
}

// AddPlayer adds a device as a player to a game. The device must be online.
func (uc *GameUseCase) AddPlayer(ctx context.Context, gameID, deviceID, nickname string, teamID *string) (*domain.Player, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game for add player: %w", err)
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}

	players, err := uc.players.ListByGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("list players for capacity check: %w", err)
	}
	if len(players) >= game.Settings.MaxPlayers {
		return nil, domain.ErrGameFull
	}

	device, err := uc.devices.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("get device %s: %w", deviceID, err)
	}
	if device.Status == domain.DeviceInGame {
		return nil, domain.ErrDeviceInGame
	}

	sessionCode, err := uc.generateUniqueSessionCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate session code: %w", err)
	}

	player := &domain.Player{
		ID:          uuid.New().String(),
		GameID:      gameID,
		TeamID:      teamID,
		DeviceID:    deviceID,
		Nickname:    nickname,
		IsAlive:     true,
		SessionCode: sessionCode,
	}
	if err := uc.players.Create(ctx, player); err != nil {
		return nil, fmt.Errorf("create player: %w", err)
	}
	return player, nil
}

// RemovePlayer removes a player and releases their device back to online.
func (uc *GameUseCase) RemovePlayer(ctx context.Context, gameID, playerID string) error {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game for remove player: %w", err)
	}
	if game.Status != domain.GameLobby {
		return domain.ErrInvalidGameState
	}

	player, err := uc.players.GetByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("get player %s: %w", playerID, err)
	}

	if err := uc.players.Delete(ctx, playerID); err != nil {
		return fmt.Errorf("delete player %s: %w", playerID, err)
	}

	// Release device back to online
	if err := uc.devices.UpdateStatus(ctx, player.DeviceID, domain.DeviceOnline); err != nil {
		slog.Error("failed to release device after player removal",
			"device_id", player.DeviceID, "player_id", playerID, "error", err)
	}
	return nil
}

func (uc *GameUseCase) ListPlayers(ctx context.Context, gameID string) ([]domain.Player, error) {
	return uc.players.ListByGame(ctx, gameID)
}

// StartGame transitions game from lobby to running within a transaction.
// Returns the game and the list of players for MQTT notification.
func (uc *GameUseCase) StartGame(ctx context.Context, gameID string) (*domain.Game, []domain.Player, error) {
	var result *domain.Game
	var resultPlayers []domain.Player

	err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		game, err := uc.games.GetByID(txCtx, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}
		if game.Status != domain.GameLobby {
			return domain.ErrInvalidGameState
		}

		now := time.Now()
		game.Status = domain.GameRunning
		game.StartedAt = &now

		if err := uc.games.Update(txCtx, game); err != nil {
			return fmt.Errorf("update game status: %w", err)
		}

		players, err := uc.players.ListByGame(txCtx, gameID)
		if err != nil {
			return fmt.Errorf("list players: %w", err)
		}

		for _, p := range players {
			if err := uc.devices.UpdateStatus(txCtx, p.DeviceID, domain.DeviceInGame); err != nil {
				return fmt.Errorf("mark device %s in-game: %w", p.DeviceID, err)
			}
		}

		result = game
		resultPlayers = players
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return result, resultPlayers, nil
}

// EndGame transitions game from running to finished within a transaction.
// Returns the game and device IDs for MQTT notification.
func (uc *GameUseCase) EndGame(ctx context.Context, gameID string) (*domain.Game, []string, error) {
	var result *domain.Game
	var deviceIDs []string

	err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		game, err := uc.games.GetByID(txCtx, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}
		if game.Status != domain.GameRunning {
			return domain.ErrInvalidGameState
		}

		now := time.Now()
		game.Status = domain.GameFinished
		game.EndedAt = &now

		if err := uc.games.Update(txCtx, game); err != nil {
			return fmt.Errorf("update game status: %w", err)
		}

		players, err := uc.players.ListByGame(txCtx, gameID)
		if err != nil {
			return fmt.Errorf("list players: %w", err)
		}

		for _, p := range players {
			if err := uc.devices.UpdateStatus(txCtx, p.DeviceID, domain.DeviceOnline); err != nil {
				return fmt.Errorf("release device %s: %w", p.DeviceID, err)
			}
			deviceIDs = append(deviceIDs, p.DeviceID)
		}

		result = game
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return result, deviceIDs, nil
}

func (uc *GameUseCase) UpdateSettings(ctx context.Context, gameID string, settings domain.GameSettings) (*domain.Game, error) {
	if err := settings.Validate(); err != nil {
		return nil, fmt.Errorf("validate settings: %w", err)
	}

	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game for update settings: %w", err)
	}
	if game.Status != domain.GameLobby {
		return nil, domain.ErrInvalidGameState
	}
	game.Settings = settings
	if err := uc.games.Update(ctx, game); err != nil {
		return nil, fmt.Errorf("update game settings: %w", err)
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

// generateUniqueSessionCode generates a unique 6-char session PIN with retry.
func (uc *GameUseCase) generateUniqueSessionCode(ctx context.Context) (string, error) {
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		code := generateGameCode()
		_, err := uc.players.GetBySessionCode(ctx, code)
		if errors.Is(err, domain.ErrNotFound) {
			return code, nil // code is unique
		}
		if err != nil {
			return "", fmt.Errorf("check session code uniqueness: %w", err)
		}
		// code exists, retry
	}
	return "", fmt.Errorf("failed to generate unique session code after %d retries", maxRetries)
}

// GetPlayerSession returns a player's session info by their PIN code (for mobile app).
func (uc *GameUseCase) GetPlayerSession(ctx context.Context, code string) (*domain.PlayerSession, error) {
	player, err := uc.players.GetBySessionCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get player by session code: %w", err)
	}

	game, err := uc.games.GetByID(ctx, player.GameID)
	if err != nil {
		return nil, fmt.Errorf("get game for session: %w", err)
	}

	session := &domain.PlayerSession{
		Player:        *player,
		Game:          *game,
		RemainingTime: uc.RemainingTime(game),
	}

	if player.TeamID != nil {
		team, err := uc.teams.GetByID(ctx, *player.TeamID)
		if err == nil {
			session.Team = team
		}
	}

	// Include leaderboard (all players sorted by score)
	leaderboard, err := uc.players.ListByGame(ctx, player.GameID)
	if err == nil {
		session.Leaderboard = leaderboard
	}

	// Include game events for kill feed
	events, err := uc.events.ListByGame(ctx, player.GameID)
	if err == nil {
		session.Events = events
	}

	return session, nil
}

// GetGameFull returns a game with its teams, players, and events.
func (uc *GameUseCase) GetGameFull(ctx context.Context, gameID string) (*domain.GameFull, error) {
	game, err := uc.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get game: %w", err)
	}
	teams, err := uc.teams.ListByGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	players, err := uc.players.ListByGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	events, err := uc.events.ListByGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	return &domain.GameFull{
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
		return fmt.Errorf("get player: %w", err)
	}
	game, err := uc.games.GetByID(ctx, player.GameID)
	if err != nil {
		return fmt.Errorf("get game: %w", err)
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
		if errors.Is(err, domain.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("get game for auto-end check: %w", err)
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
