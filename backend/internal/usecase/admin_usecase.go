package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/config"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase struct {
	userRepo    domain.UserRepository
	sessionRepo domain.AdminSessionRepository
	gameUC      *GameUseCase
	playerRepo  domain.GamePlayerRepository
	cache       domain.GameCache
	eventBus    domain.EventBus
	cfg         config.JWTConfig
}

func NewAdminUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.AdminSessionRepository,
	gameUC *GameUseCase,
	playerRepo domain.GamePlayerRepository,
	cache domain.GameCache,
	eventBus domain.EventBus,
	cfg config.JWTConfig,
) *AdminUseCase {
	return &AdminUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		gameUC:      gameUC,
		playerRepo:  playerRepo,
		cache:       cache,
		eventBus:    eventBus,
		cfg:         cfg,
	}
}

func (uc *AdminUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Delete old sessions
	_ = uc.sessionRepo.DeleteByUserID(ctx, user.ID)

	// Create new session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	session := &domain.AdminSession{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(uc.cfg.Expiration),
		CreatedAt: time.Now(),
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}
	return token, nil
}

func (uc *AdminUseCase) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	session, err := uc.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrSessionExpired
	}
	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return user, nil
}

func (uc *AdminUseCase) ExecuteControl(ctx context.Context, cmd domain.AdminControlCommand) error {
	switch cmd.Action {
	case domain.AdminActionPause:
		return uc.gameUC.PauseGame(ctx, cmd.GameID)
	case domain.AdminActionResume:
		return uc.gameUC.ResumeGame(ctx, cmd.GameID)
	case domain.AdminActionEnd:
		return uc.gameUC.EndGame(ctx, cmd.GameID)
	case domain.AdminActionRestart:
		return uc.restartGame(ctx, cmd.GameID)
	case domain.AdminActionRevive:
		return uc.revivePlayer(ctx, cmd.GameID, cmd.PlayerID)
	case domain.AdminActionKick:
		return uc.kickPlayer(ctx, cmd.GameID, cmd.PlayerID)
	case domain.AdminActionChangeTeam:
		teamID, _ := cmd.Params["team_id"].(string)
		return uc.changeTeam(ctx, cmd.GameID, cmd.PlayerID, teamID)
	default:
		return domain.ErrInvalidAction
	}
}

func (uc *AdminUseCase) restartGame(ctx context.Context, gameID string) error {
	game, err := uc.gameUC.GetGame(ctx, gameID)
	if err != nil {
		return err
	}
	// End current game first
	_ = uc.gameUC.EndGame(ctx, gameID)

	// Reset game state
	game.Status = domain.GameStatusPending
	game.StartedAt = nil
	game.EndedAt = nil
	if err := uc.gameUC.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	// Reset players
	players, _ := uc.playerRepo.GetByGameID(ctx, gameID)
	for _, p := range players {
		p.HP = game.Config.Player.MaxHP
		p.Score = 0
		p.Kills = 0
		p.Deaths = 0
		p.IsAlive = true
		p.LivesRemaining = game.Config.Player.Lives
		_ = uc.playerRepo.Update(ctx, p)
	}
	return nil
}

func (uc *AdminUseCase) revivePlayer(ctx context.Context, gameID, playerID string) error {
	game, err := uc.gameUC.GetGame(ctx, gameID)
	if err != nil {
		return err
	}
	state, err := uc.cache.GetPlayerState(ctx, gameID, playerID)
	if err != nil {
		return domain.ErrPlayerNotInGame
	}
	state.IsAlive = true
	state.HP = game.Config.Player.MaxHP
	_ = uc.cache.SetPlayerState(ctx, gameID, state)

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypePlayerRevive),
		GameID: gameID,
		Payload: map[string]interface{}{
			"player_id": playerID,
			"hp":        game.Config.Player.MaxHP,
		},
	})
	return nil
}

func (uc *AdminUseCase) kickPlayer(ctx context.Context, gameID, playerID string) error {
	state, err := uc.cache.GetPlayerState(ctx, gameID, playerID)
	if err != nil {
		return domain.ErrPlayerNotInGame
	}
	state.IsAlive = false
	state.HP = 0
	state.LivesRemaining = 0
	_ = uc.cache.SetPlayerState(ctx, gameID, state)

	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypePlayerKick),
		GameID: gameID,
		Payload: map[string]interface{}{
			"player_id": playerID,
		},
	})
	return nil
}

func (uc *AdminUseCase) changeTeam(ctx context.Context, gameID, playerID, teamID string) error {
	player, err := uc.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return domain.ErrPlayerNotInGame
	}
	player.TeamID = &teamID
	if err := uc.playerRepo.Update(ctx, player); err != nil {
		return err
	}
	uc.eventBus.Publish(gameID, domain.WSMessage{
		Type:   string(domain.EventTypeTeamChange),
		GameID: gameID,
		Payload: map[string]interface{}{
			"player_id": playerID,
			"team_id":   teamID,
		},
	})
	return nil
}
