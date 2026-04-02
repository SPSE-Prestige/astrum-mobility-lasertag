package domain

import (
	"context"
	"time"
)

// ── Repository interfaces ──

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, s *Session) error
	GetByToken(ctx context.Context, token string) (*Session, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

type DeviceRepository interface {
	Upsert(ctx context.Context, d *Device) error
	GetByDeviceID(ctx context.Context, deviceID string) (*Device, error)
	ListAll(ctx context.Context) ([]Device, error)
	ListByStatus(ctx context.Context, status DeviceStatus) ([]Device, error)
	UpdateStatus(ctx context.Context, deviceID string, status DeviceStatus) error
	UpdateLastSeen(ctx context.Context, deviceID string) error
}

type GameRepository interface {
	Create(ctx context.Context, g *Game) error
	GetByID(ctx context.Context, id string) (*Game, error)
	GetByCode(ctx context.Context, code string) (*Game, error)
	Update(ctx context.Context, g *Game) error
	ListAll(ctx context.Context) ([]Game, error)
	ListByStatus(ctx context.Context, status GameStatus) ([]Game, error)
}

type TeamRepository interface {
	Create(ctx context.Context, t *Team) error
	GetByID(ctx context.Context, id int) (*Team, error)
	ListByGame(ctx context.Context, gameID string) ([]Team, error)
	Delete(ctx context.Context, id int) error
}

type PlayerRepository interface {
	Create(ctx context.Context, p *Player) error
	GetByID(ctx context.Context, id string) (*Player, error)
	GetByGameAndDevice(ctx context.Context, gameID, deviceID string) (*Player, error)
	// FindActivePlayerByDevice finds a player with this device in a running game (if any).
	FindActivePlayerByDevice(ctx context.Context, deviceID string) (*Player, *Game, error)
	ListByGame(ctx context.Context, gameID string) ([]Player, error)
	ListByTeam(ctx context.Context, teamID int) ([]Player, error)
	Update(ctx context.Context, p *Player) error
	Delete(ctx context.Context, id string) error
	// KillPlayer atomically sets is_alive=false, increments deaths, resets kill_streak and weapon_level.
	// Returns false if already dead.
	KillPlayer(ctx context.Context, playerID string) (bool, error)
	// AddKillScore atomically increments kills, score, and kill_streak.
	// If killsPerUpgrade > 0 and streak hits the threshold, weapon_level is incremented.
	// Returns the updated streak state.
	AddKillScore(ctx context.Context, playerID string, score, killsPerUpgrade int) (*KillScoreResult, error)
	// SubKillScore atomically decrements kills and score, resets streak (friendly fire penalty).
	SubKillScore(ctx context.Context, playerID string, score int) error
	// Respawn atomically sets is_alive=true.
	Respawn(ctx context.Context, playerID string) error
	// IncrementShotsFired atomically increments the shots_fired counter.
	IncrementShotsFired(ctx context.Context, playerID string) error
	// GetBySessionCode finds a player by their unique session PIN code.
	GetBySessionCode(ctx context.Context, code string) (*Player, error)
}

type EventRepository interface {
	Create(ctx context.Context, e *GameEvent) error
	ListByGame(ctx context.Context, gameID string) ([]GameEvent, error)
}

// ── Transaction support ──

// KillScoreResult holds the streak state returned after a kill score update.
type KillScoreResult struct {
	KillStreak  int
	WeaponLevel int
}

// TxManager provides transactional execution.
type TxManager interface {
	// WithTx executes fn within a database transaction.
	// If fn returns an error, the transaction is rolled back; otherwise it is committed.
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// ── Use case port interfaces ──

// AuthUseCase defines authentication operations.
type AuthUseCasePort interface {
	Login(ctx context.Context, username, password string) (*Session, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	Logout(ctx context.Context, token string) error
	CleanupExpiredSessions(ctx context.Context) error
}

// DeviceUseCasePort defines device management operations.
type DeviceUseCasePort interface {
	Register(ctx context.Context, deviceID string) (*Device, error)
	Heartbeat(ctx context.Context, deviceID string) error
	// Reconnect checks if a device has an active game session and restores its state.
	// Returns nil info if the device is not in any running game.
	Reconnect(ctx context.Context, deviceId string) (*ReconnectInfo, error)
	MarkOffline(ctx context.Context, timeout time.Duration) ([]string, error)
	ListAll(ctx context.Context) ([]Device, error)
	ListAvailable(ctx context.Context) ([]Device, error)
}

// GameUseCasePort defines game management operations.
type GameUseCasePort interface {
	CreateGame(ctx context.Context, settings *GameSettings) (*Game, error)
	GetGame(ctx context.Context, id string) (*Game, error)
	ListGames(ctx context.Context) ([]Game, error)
	GetGameFull(ctx context.Context, gameID string) (*GameFull, error)
	StartGame(ctx context.Context, gameID string) (*Game, []Player, error) // returns game + device IDs
	EndGame(ctx context.Context, gameID string) (*Game, []string, error)   // returns game + device IDs
	UpdateSettings(ctx context.Context, gameID string, settings GameSettings) (*Game, error)
	AddTeam(ctx context.Context, gameID, name, color string) (*Team, error)
	ListTeams(ctx context.Context, gameID string) ([]Team, error)
	RemoveTeam(ctx context.Context, gameID string, teamID int) error
	AddPlayer(ctx context.Context, gameID, deviceID, nickname string, teamID *int) (*Player, error)
	RemovePlayer(ctx context.Context, gameID, playerID string) error
	ListPlayers(ctx context.Context, gameID string) ([]Player, error)
	GetLeaderboard(ctx context.Context, gameID string) ([]Player, error)
	ListEvents(ctx context.Context, gameID string) ([]GameEvent, error)
	UpdatePlayerTeam(ctx context.Context, playerID string, teamID *int) error
	ShouldAutoEnd(ctx context.Context, gameID string) (bool, error)
	RemainingTime(game *Game) int
	// GetPlayerSession returns a player's game session by their session PIN code.
	GetPlayerSession(ctx context.Context, code string) (*PlayerSession, error)
}

// HitUseCasePort defines hit processing operations.
type HitUseCasePort interface {
	ProcessHit(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*HitResult, error)
	Respawn(ctx context.Context, gameID, deviceID string) error
	// RecordShot increments the shots_fired counter for a device in a running game.
	RecordShot(ctx context.Context, gameID, deviceID string) error
}

// ── Infrastructure port interfaces ──

// MQTTPublisher defines MQTT command publishing.
type MQTTPublisher interface {
	SendCommand(deviceID string, command any)
	PublishGameStart(players []Player, game Game)
	PublishGameEnd(deviceIDs []string)
	// PublishGameState sends current game state to a single reconnecting device.
	PublishGameState(deviceID string, info *ReconnectInfo)
}

// WSBroadcaster defines WebSocket broadcasting.
type WSBroadcaster interface {
	Broadcast(msg any)
}
