package mocks

import (
	"context"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

// Ensure unused imports are consumed.
var _ time.Duration

// ── UserRepository ──

type MockUserRepository struct {
	GetByIDFn       func(ctx context.Context, id string) (*domain.User, error)
	GetByUsernameFn func(ctx context.Context, username string) (*domain.User, error)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return m.GetByUsernameFn(ctx, username)
}

// ── SessionRepository ──

type MockSessionRepository struct {
	CreateFn        func(ctx context.Context, s *domain.Session) error
	GetByTokenFn    func(ctx context.Context, token string) (*domain.Session, error)
	DeleteByTokenFn func(ctx context.Context, token string) error
	DeleteExpiredFn func(ctx context.Context) error
}

func (m *MockSessionRepository) Create(ctx context.Context, s *domain.Session) error {
	return m.CreateFn(ctx, s)
}

func (m *MockSessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	return m.GetByTokenFn(ctx, token)
}

func (m *MockSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return m.DeleteByTokenFn(ctx, token)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	return m.DeleteExpiredFn(ctx)
}

// ── DeviceRepository ──

type MockDeviceRepository struct {
	UpsertFn         func(ctx context.Context, d *domain.Device) error
	GetByDeviceIDFn  func(ctx context.Context, deviceID string) (*domain.Device, error)
	ListAllFn        func(ctx context.Context) ([]domain.Device, error)
	ListByStatusFn   func(ctx context.Context, status domain.DeviceStatus) ([]domain.Device, error)
	UpdateStatusFn   func(ctx context.Context, deviceID string, status domain.DeviceStatus) error
	UpdateLastSeenFn func(ctx context.Context, deviceID string) error
}

func (m *MockDeviceRepository) Upsert(ctx context.Context, d *domain.Device) error {
	return m.UpsertFn(ctx, d)
}

func (m *MockDeviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	return m.GetByDeviceIDFn(ctx, deviceID)
}

func (m *MockDeviceRepository) ListAll(ctx context.Context) ([]domain.Device, error) {
	return m.ListAllFn(ctx)
}

func (m *MockDeviceRepository) ListByStatus(ctx context.Context, status domain.DeviceStatus) ([]domain.Device, error) {
	return m.ListByStatusFn(ctx, status)
}

func (m *MockDeviceRepository) UpdateStatus(ctx context.Context, deviceID string, status domain.DeviceStatus) error {
	return m.UpdateStatusFn(ctx, deviceID, status)
}

func (m *MockDeviceRepository) UpdateLastSeen(ctx context.Context, deviceID string) error {
	return m.UpdateLastSeenFn(ctx, deviceID)
}

// ── GameRepository ──

type MockGameRepository struct {
	CreateFn       func(ctx context.Context, g *domain.Game) error
	GetByIDFn      func(ctx context.Context, id string) (*domain.Game, error)
	GetByCodeFn    func(ctx context.Context, code string) (*domain.Game, error)
	UpdateFn       func(ctx context.Context, g *domain.Game) error
	ListAllFn      func(ctx context.Context) ([]domain.Game, error)
	ListByStatusFn func(ctx context.Context, status domain.GameStatus) ([]domain.Game, error)
}

func (m *MockGameRepository) Create(ctx context.Context, g *domain.Game) error {
	return m.CreateFn(ctx, g)
}

func (m *MockGameRepository) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockGameRepository) GetByCode(ctx context.Context, code string) (*domain.Game, error) {
	return m.GetByCodeFn(ctx, code)
}

func (m *MockGameRepository) Update(ctx context.Context, g *domain.Game) error {
	return m.UpdateFn(ctx, g)
}

func (m *MockGameRepository) ListAll(ctx context.Context) ([]domain.Game, error) {
	return m.ListAllFn(ctx)
}

func (m *MockGameRepository) ListByStatus(ctx context.Context, status domain.GameStatus) ([]domain.Game, error) {
	return m.ListByStatusFn(ctx, status)
}

// ── TeamRepository ──

type MockTeamRepository struct {
	CreateFn     func(ctx context.Context, t *domain.Team) error
	GetByIDFn    func(ctx context.Context, id string) (*domain.Team, error)
	ListByGameFn func(ctx context.Context, gameID string) ([]domain.Team, error)
	DeleteFn     func(ctx context.Context, id string) error
}

func (m *MockTeamRepository) Create(ctx context.Context, t *domain.Team) error {
	return m.CreateFn(ctx, t)
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockTeamRepository) ListByGame(ctx context.Context, gameID string) ([]domain.Team, error) {
	return m.ListByGameFn(ctx, gameID)
}

func (m *MockTeamRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

// ── PlayerRepository ──

type MockPlayerRepository struct {
	CreateFn             func(ctx context.Context, p *domain.Player) error
	GetByIDFn            func(ctx context.Context, id string) (*domain.Player, error)
	GetByGameAndDeviceFn func(ctx context.Context, gameID, deviceID string) (*domain.Player, error)
	ListByGameFn         func(ctx context.Context, gameID string) ([]domain.Player, error)
	ListByTeamFn         func(ctx context.Context, teamID string) ([]domain.Player, error)
	UpdateFn             func(ctx context.Context, p *domain.Player) error
	DeleteFn             func(ctx context.Context, id string) error
	KillPlayerFn         func(ctx context.Context, playerID string) (bool, error)
	AddKillScoreFn       func(ctx context.Context, playerID string, score, killsPerUpgrade int) (*domain.KillScoreResult, error)
	RespawnFn            func(ctx context.Context, playerID string) error
}

func (m *MockPlayerRepository) Create(ctx context.Context, p *domain.Player) error {
	return m.CreateFn(ctx, p)
}

func (m *MockPlayerRepository) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockPlayerRepository) GetByGameAndDevice(ctx context.Context, gameID, deviceID string) (*domain.Player, error) {
	return m.GetByGameAndDeviceFn(ctx, gameID, deviceID)
}

func (m *MockPlayerRepository) ListByGame(ctx context.Context, gameID string) ([]domain.Player, error) {
	return m.ListByGameFn(ctx, gameID)
}

func (m *MockPlayerRepository) ListByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	return m.ListByTeamFn(ctx, teamID)
}

func (m *MockPlayerRepository) Update(ctx context.Context, p *domain.Player) error {
	return m.UpdateFn(ctx, p)
}

func (m *MockPlayerRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockPlayerRepository) KillPlayer(ctx context.Context, playerID string) (bool, error) {
	return m.KillPlayerFn(ctx, playerID)
}

func (m *MockPlayerRepository) AddKillScore(ctx context.Context, playerID string, score, killsPerUpgrade int) (*domain.KillScoreResult, error) {
	return m.AddKillScoreFn(ctx, playerID, score, killsPerUpgrade)
}

func (m *MockPlayerRepository) Respawn(ctx context.Context, playerID string) error {
	return m.RespawnFn(ctx, playerID)
}

// ── EventRepository ──

type MockEventRepository struct {
	CreateFn     func(ctx context.Context, e *domain.GameEvent) error
	ListByGameFn func(ctx context.Context, gameID string) ([]domain.GameEvent, error)
}

func (m *MockEventRepository) Create(ctx context.Context, e *domain.GameEvent) error {
	return m.CreateFn(ctx, e)
}

func (m *MockEventRepository) ListByGame(ctx context.Context, gameID string) ([]domain.GameEvent, error) {
	return m.ListByGameFn(ctx, gameID)
}

// ── TxManager ──

type MockTxManager struct {
	WithTxFn func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.WithTxFn != nil {
		return m.WithTxFn(ctx, fn)
	}
	// Default: execute fn directly without a real transaction.
	return fn(ctx)
}

// ── MQTTPublisher ──

type MockMQTTPublisher struct {
	SendCommandFn      func(deviceID string, command any)
	PublishGameStartFn func(deviceIDs []string, gameID string)
	PublishGameEndFn   func(deviceIDs []string)
}

func (m *MockMQTTPublisher) SendCommand(deviceID string, command any) {
	if m.SendCommandFn != nil {
		m.SendCommandFn(deviceID, command)
	}
}

func (m *MockMQTTPublisher) PublishGameStart(deviceIDs []string, gameID string) {
	if m.PublishGameStartFn != nil {
		m.PublishGameStartFn(deviceIDs, gameID)
	}
}

func (m *MockMQTTPublisher) PublishGameEnd(deviceIDs []string) {
	if m.PublishGameEndFn != nil {
		m.PublishGameEndFn(deviceIDs)
	}
}

// ── WSBroadcaster ──

type MockWSBroadcaster struct {
	BroadcastFn func(msg any)
}

func (m *MockWSBroadcaster) Broadcast(msg any) {
	if m.BroadcastFn != nil {
		m.BroadcastFn(msg)
	}
}

// ── AuthUseCasePort ──

type MockAuthUseCasePort struct {
	LoginFn                  func(ctx context.Context, username, password string) (*domain.Session, error)
	ValidateTokenFn          func(ctx context.Context, token string) (*domain.User, error)
	LogoutFn                 func(ctx context.Context, token string) error
	CleanupExpiredSessionsFn func(ctx context.Context) error
}

func (m *MockAuthUseCasePort) Login(ctx context.Context, username, password string) (*domain.Session, error) {
	return m.LoginFn(ctx, username, password)
}

func (m *MockAuthUseCasePort) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	return m.ValidateTokenFn(ctx, token)
}

func (m *MockAuthUseCasePort) Logout(ctx context.Context, token string) error {
	return m.LogoutFn(ctx, token)
}

func (m *MockAuthUseCasePort) CleanupExpiredSessions(ctx context.Context) error {
	return m.CleanupExpiredSessionsFn(ctx)
}

// ── DeviceUseCasePort ──

type MockDeviceUseCasePort struct {
	RegisterFn      func(ctx context.Context, deviceID string) (*domain.Device, error)
	HeartbeatFn     func(ctx context.Context, deviceID string) error
	MarkOfflineFn   func(ctx context.Context, timeout time.Duration) ([]string, error)
	ListAllFn       func(ctx context.Context) ([]domain.Device, error)
	ListAvailableFn func(ctx context.Context) ([]domain.Device, error)
}

func (m *MockDeviceUseCasePort) Register(ctx context.Context, deviceID string) (*domain.Device, error) {
	return m.RegisterFn(ctx, deviceID)
}

func (m *MockDeviceUseCasePort) Heartbeat(ctx context.Context, deviceID string) error {
	return m.HeartbeatFn(ctx, deviceID)
}

func (m *MockDeviceUseCasePort) MarkOffline(ctx context.Context, timeout time.Duration) ([]string, error) {
	return m.MarkOfflineFn(ctx, timeout)
}

func (m *MockDeviceUseCasePort) ListAll(ctx context.Context) ([]domain.Device, error) {
	return m.ListAllFn(ctx)
}

func (m *MockDeviceUseCasePort) ListAvailable(ctx context.Context) ([]domain.Device, error) {
	return m.ListAvailableFn(ctx)
}

// ── GameUseCasePort ──

type MockGameUseCasePort struct {
	CreateGameFn      func(ctx context.Context, settings *domain.GameSettings) (*domain.Game, error)
	GetGameFn         func(ctx context.Context, id string) (*domain.Game, error)
	ListGamesFn       func(ctx context.Context) ([]domain.Game, error)
	GetGameFullFn     func(ctx context.Context, gameID string) (*domain.GameFull, error)
	StartGameFn       func(ctx context.Context, gameID string) (*domain.Game, []string, error)
	EndGameFn         func(ctx context.Context, gameID string) (*domain.Game, []string, error)
	UpdateSettingsFn  func(ctx context.Context, gameID string, settings domain.GameSettings) (*domain.Game, error)
	AddTeamFn         func(ctx context.Context, gameID, name, color string) (*domain.Team, error)
	ListTeamsFn       func(ctx context.Context, gameID string) ([]domain.Team, error)
	RemoveTeamFn      func(ctx context.Context, gameID, teamID string) error
	AddPlayerFn       func(ctx context.Context, gameID, deviceID, nickname string, teamID *string) (*domain.Player, error)
	RemovePlayerFn    func(ctx context.Context, gameID, playerID string) error
	ListPlayersFn     func(ctx context.Context, gameID string) ([]domain.Player, error)
	GetLeaderboardFn  func(ctx context.Context, gameID string) ([]domain.Player, error)
	ListEventsFn      func(ctx context.Context, gameID string) ([]domain.GameEvent, error)
	UpdatePlayerTeamFn func(ctx context.Context, playerID string, teamID *string) error
	ShouldAutoEndFn   func(ctx context.Context, gameID string) (bool, error)
	RemainingTimeFn   func(game *domain.Game) int
}

func (m *MockGameUseCasePort) CreateGame(ctx context.Context, settings *domain.GameSettings) (*domain.Game, error) {
	return m.CreateGameFn(ctx, settings)
}

func (m *MockGameUseCasePort) GetGame(ctx context.Context, id string) (*domain.Game, error) {
	return m.GetGameFn(ctx, id)
}

func (m *MockGameUseCasePort) ListGames(ctx context.Context) ([]domain.Game, error) {
	return m.ListGamesFn(ctx)
}

func (m *MockGameUseCasePort) GetGameFull(ctx context.Context, gameID string) (*domain.GameFull, error) {
	return m.GetGameFullFn(ctx, gameID)
}

func (m *MockGameUseCasePort) StartGame(ctx context.Context, gameID string) (*domain.Game, []string, error) {
	return m.StartGameFn(ctx, gameID)
}

func (m *MockGameUseCasePort) EndGame(ctx context.Context, gameID string) (*domain.Game, []string, error) {
	return m.EndGameFn(ctx, gameID)
}

func (m *MockGameUseCasePort) UpdateSettings(ctx context.Context, gameID string, settings domain.GameSettings) (*domain.Game, error) {
	return m.UpdateSettingsFn(ctx, gameID, settings)
}

func (m *MockGameUseCasePort) AddTeam(ctx context.Context, gameID, name, color string) (*domain.Team, error) {
	return m.AddTeamFn(ctx, gameID, name, color)
}

func (m *MockGameUseCasePort) ListTeams(ctx context.Context, gameID string) ([]domain.Team, error) {
	return m.ListTeamsFn(ctx, gameID)
}

func (m *MockGameUseCasePort) RemoveTeam(ctx context.Context, gameID, teamID string) error {
	return m.RemoveTeamFn(ctx, gameID, teamID)
}

func (m *MockGameUseCasePort) AddPlayer(ctx context.Context, gameID, deviceID, nickname string, teamID *string) (*domain.Player, error) {
	return m.AddPlayerFn(ctx, gameID, deviceID, nickname, teamID)
}

func (m *MockGameUseCasePort) RemovePlayer(ctx context.Context, gameID, playerID string) error {
	return m.RemovePlayerFn(ctx, gameID, playerID)
}

func (m *MockGameUseCasePort) ListPlayers(ctx context.Context, gameID string) ([]domain.Player, error) {
	return m.ListPlayersFn(ctx, gameID)
}

func (m *MockGameUseCasePort) GetLeaderboard(ctx context.Context, gameID string) ([]domain.Player, error) {
	return m.GetLeaderboardFn(ctx, gameID)
}

func (m *MockGameUseCasePort) ListEvents(ctx context.Context, gameID string) ([]domain.GameEvent, error) {
	return m.ListEventsFn(ctx, gameID)
}

func (m *MockGameUseCasePort) UpdatePlayerTeam(ctx context.Context, playerID string, teamID *string) error {
	return m.UpdatePlayerTeamFn(ctx, playerID, teamID)
}

func (m *MockGameUseCasePort) ShouldAutoEnd(ctx context.Context, gameID string) (bool, error) {
	return m.ShouldAutoEndFn(ctx, gameID)
}

func (m *MockGameUseCasePort) RemainingTime(game *domain.Game) int {
	return m.RemainingTimeFn(game)
}

// ── HitUseCasePort ──

type MockHitUseCasePort struct {
	ProcessHitFn func(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*domain.HitResult, error)
	RespawnFn    func(ctx context.Context, gameID, deviceID string) error
}

func (m *MockHitUseCasePort) ProcessHit(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*domain.HitResult, error) {
	return m.ProcessHitFn(ctx, gameID, attackerDeviceID, victimDeviceID)
}

func (m *MockHitUseCasePort) Respawn(ctx context.Context, gameID, deviceID string) error {
	return m.RespawnFn(ctx, gameID, deviceID)
}
