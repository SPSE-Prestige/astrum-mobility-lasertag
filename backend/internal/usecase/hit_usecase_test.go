package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain/mocks"
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int { return &i }

func newRunningGame() *domain.Game {
	return &domain.Game{
		ID:     "game-1",
		Code:   "ABCD",
		Status: domain.GameRunning,
		Settings: domain.GameSettings{
			FriendlyFire: false,
			ScorePerKill: 100,
		},
	}
}

func newHitUseCase(
	games *mocks.MockGameRepository,
	players *mocks.MockPlayerRepository,
	events *mocks.MockEventRepository,
	txMgr *mocks.MockTxManager,
) *HitUseCase {
	return NewHitUseCase(games, players, events, txMgr)
}

func TestProcessHit_Success(t *testing.T) {
	game := newRunningGame()
	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		TeamID: ptrInt(1), Nickname: "Alice", IsAlive: true,
		Score: 0, Kills: 0,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		TeamID: ptrInt(2), Nickname: "Bob", IsAlive: true,
		Score: 0, Kills: 0,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, id string) (*domain.Game, error) {
			if id == game.ID {
				return game, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, gameID, deviceID string) (*domain.Player, error) {
			switch deviceID {
			case "dev-A":
				return attacker, nil
			case "dev-B":
				return victim, nil
			}
			return nil, domain.ErrNotFound
		},
		KillPlayerFn: func(_ context.Context, playerID string) (bool, error) {
			if playerID == victim.ID {
				return true, nil
			}
			return false, nil
		},
		AddKillScoreFn: func(_ context.Context, playerID string, score, killsPerUpgrade int) (*domain.KillScoreResult, error) {
			if playerID != attacker.ID {
				t.Errorf("expected attacker ID %q, got %q", attacker.ID, playerID)
			}
			if score != game.Settings.ScorePerKill {
				t.Errorf("expected score %d, got %d", game.Settings.ScorePerKill, score)
			}
			return &domain.KillScoreResult{KillStreak: 1, WeaponLevel: 0}, nil
		},
	}
	events := &mocks.MockEventRepository{
		CreateFn: func(_ context.Context, e *domain.GameEvent) error {
			if e.Type != "kill" {
				t.Errorf("expected event type 'kill', got %q", e.Type)
			}
			return nil
		},
	}
	txMgr := &mocks.MockTxManager{
		WithTxFn: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	}

	uc := newHitUseCase(games, players, events, txMgr)
	result, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Kill {
		t.Error("expected Kill to be true")
	}
	if result.AttackerID != attacker.ID {
		t.Errorf("expected AttackerID %q, got %q", attacker.ID, result.AttackerID)
	}
	if result.VictimID != victim.ID {
		t.Errorf("expected VictimID %q, got %q", victim.ID, result.VictimID)
	}
	if result.AttackerScore != 100 {
		t.Errorf("expected AttackerScore 100, got %d", result.AttackerScore)
	}
	if result.AttackerKills != 1 {
		t.Errorf("expected AttackerKills 1, got %d", result.AttackerKills)
	}
}

func TestProcessHit_GameNotRunning(t *testing.T) {
	game := &domain.Game{ID: "game-1", Status: domain.GameLobby}
	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}

	uc := newHitUseCase(games, nil, nil, nil)
	_, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if !errors.Is(err, domain.ErrInvalidGameState) {
		t.Errorf("expected ErrInvalidGameState, got %v", err)
	}
}

func TestProcessHit_SelfHit(t *testing.T) {
	game := newRunningGame()
	player := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		IsAlive: true,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, _ string) (*domain.Player, error) {
			return player, nil
		},
	}

	uc := newHitUseCase(games, players, nil, nil)
	_, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-A")
	if !errors.Is(err, domain.ErrSelfHit) {
		t.Errorf("expected ErrSelfHit, got %v", err)
	}
}

func TestProcessHit_VictimAlreadyDead(t *testing.T) {
	game := newRunningGame()
	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		TeamID: ptrInt(1), IsAlive: true,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		TeamID: ptrInt(2), IsAlive: false,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, deviceID string) (*domain.Player, error) {
			if deviceID == "dev-A" {
				return attacker, nil
			}
			return victim, nil
		},
	}

	uc := newHitUseCase(games, players, nil, nil)
	_, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if !errors.Is(err, domain.ErrPlayerDead) {
		t.Errorf("expected ErrPlayerDead, got %v", err)
	}
}

func TestProcessHit_FriendlyFire(t *testing.T) {
	game := newRunningGame()
	game.Settings.FriendlyFire = false
	sameTeam := ptrInt(1)

	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		TeamID: sameTeam, IsAlive: true,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		TeamID: sameTeam, IsAlive: true,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, deviceID string) (*domain.Player, error) {
			if deviceID == "dev-A" {
				return attacker, nil
			}
			return victim, nil
		},
	}

	uc := newHitUseCase(games, players, nil, nil)
	_, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if !errors.Is(err, domain.ErrFriendlyFire) {
		t.Errorf("expected ErrFriendlyFire, got %v", err)
	}
}

func TestProcessHit_FriendlyFireAllowed(t *testing.T) {
	game := newRunningGame()
	game.Settings.FriendlyFire = true
	sameTeam := ptrInt(1)

	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		TeamID: sameTeam, Nickname: "Alice", IsAlive: true,
		Score: 200, Kills: 2,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		TeamID: sameTeam, Nickname: "Bob", IsAlive: true,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, deviceID string) (*domain.Player, error) {
			if deviceID == "dev-A" {
				return attacker, nil
			}
			return victim, nil
		},
		KillPlayerFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
		AddKillScoreFn: func(_ context.Context, _ string, _, _ int) (*domain.KillScoreResult, error) {
			return &domain.KillScoreResult{KillStreak: 1, WeaponLevel: 0}, nil
		},
	}
	events := &mocks.MockEventRepository{
		CreateFn: func(_ context.Context, _ *domain.GameEvent) error {
			return nil
		},
	}
	txMgr := &mocks.MockTxManager{
		WithTxFn: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	}

	uc := newHitUseCase(games, players, events, txMgr)
	result, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Kill {
		t.Error("expected Kill to be true with friendly fire allowed")
	}
	if result.AttackerScore != 300 {
		t.Errorf("expected AttackerScore 300, got %d", result.AttackerScore)
	}
	if result.AttackerKills != 3 {
		t.Errorf("expected AttackerKills 3, got %d", result.AttackerKills)
	}
}

func TestProcessHit_WeaponUpgrade(t *testing.T) {
	game := newRunningGame()
	game.Settings.KillsPerUpgrade = 3

	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		Nickname: "Alice", IsAlive: true, Score: 200, Kills: 2,
		KillStreak: 2, WeaponLevel: 0,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		Nickname: "Bob", IsAlive: true,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, deviceID string) (*domain.Player, error) {
			if deviceID == "dev-A" {
				return attacker, nil
			}
			return victim, nil
		},
		KillPlayerFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
		AddKillScoreFn: func(_ context.Context, _ string, _, killsPerUpgrade int) (*domain.KillScoreResult, error) {
			if killsPerUpgrade != 3 {
				t.Errorf("expected killsPerUpgrade 3, got %d", killsPerUpgrade)
			}
			// Streak 3 with killsPerUpgrade 3 → upgrade to level 1
			return &domain.KillScoreResult{KillStreak: 3, WeaponLevel: 1}, nil
		},
	}
	events := &mocks.MockEventRepository{
		CreateFn: func(_ context.Context, _ *domain.GameEvent) error {
			return nil
		},
	}
	txMgr := &mocks.MockTxManager{
		WithTxFn: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	}

	uc := newHitUseCase(games, players, events, txMgr)
	result, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.WeaponUpgraded {
		t.Error("expected WeaponUpgraded to be true")
	}
	if result.WeaponLevel != 1 {
		t.Errorf("expected WeaponLevel 1, got %d", result.WeaponLevel)
	}
	if result.KillStreak != 3 {
		t.Errorf("expected KillStreak 3, got %d", result.KillStreak)
	}
}

func TestProcessHit_NoUpgradeWhenDisabled(t *testing.T) {
	game := newRunningGame()
	game.Settings.KillsPerUpgrade = 0 // disabled

	attacker := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		Nickname: "Alice", IsAlive: true, Score: 200, Kills: 2,
	}
	victim := &domain.Player{
		ID: "player-2", GameID: "game-1", DeviceID: "dev-B",
		Nickname: "Bob", IsAlive: true,
	}

	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, deviceID string) (*domain.Player, error) {
			if deviceID == "dev-A" {
				return attacker, nil
			}
			return victim, nil
		},
		KillPlayerFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
		AddKillScoreFn: func(_ context.Context, _ string, _, _ int) (*domain.KillScoreResult, error) {
			return &domain.KillScoreResult{KillStreak: 3, WeaponLevel: 0}, nil
		},
	}
	events := &mocks.MockEventRepository{
		CreateFn: func(_ context.Context, _ *domain.GameEvent) error {
			return nil
		},
	}
	txMgr := &mocks.MockTxManager{
		WithTxFn: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	}

	uc := newHitUseCase(games, players, events, txMgr)
	result, err := uc.ProcessHit(context.Background(), "game-1", "dev-A", "dev-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WeaponUpgraded {
		t.Error("expected WeaponUpgraded to be false when killsPerUpgrade=0")
	}
	if result.WeaponLevel != 0 {
		t.Errorf("expected WeaponLevel 0, got %d", result.WeaponLevel)
	}
}

func TestRespawn_Success(t *testing.T) {
	game := newRunningGame()
	player := &domain.Player{
		ID: "player-1", GameID: "game-1", DeviceID: "dev-A",
		Nickname: "Alice", IsAlive: false,
	}

	respawnCalled := false
	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}
	players := &mocks.MockPlayerRepository{
		GetByGameAndDeviceFn: func(_ context.Context, _, _ string) (*domain.Player, error) {
			return player, nil
		},
		RespawnFn: func(_ context.Context, playerID string) error {
			respawnCalled = true
			if playerID != player.ID {
				t.Errorf("expected player ID %q, got %q", player.ID, playerID)
			}
			return nil
		},
	}
	events := &mocks.MockEventRepository{
		CreateFn: func(_ context.Context, e *domain.GameEvent) error {
			if e.Type != "respawn" {
				t.Errorf("expected event type 'respawn', got %q", e.Type)
			}
			return nil
		},
	}

	uc := newHitUseCase(games, players, events, nil)
	err := uc.Respawn(context.Background(), "game-1", "dev-A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !respawnCalled {
		t.Error("expected Respawn to be called on player repository")
	}
}

func TestRespawn_GameNotRunning(t *testing.T) {
	game := &domain.Game{ID: "game-1", Status: domain.GameFinished}
	games := &mocks.MockGameRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.Game, error) {
			return game, nil
		},
	}

	uc := newHitUseCase(games, nil, nil, nil)
	err := uc.Respawn(context.Background(), "game-1", "dev-A")
	if !errors.Is(err, domain.ErrInvalidGameState) {
		t.Errorf("expected ErrInvalidGameState, got %v", err)
	}
}
