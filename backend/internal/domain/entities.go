package domain

import (
	"fmt"
	"time"
)

// ── Enums ──

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type DeviceStatus string

const (
	DeviceOnline  DeviceStatus = "online"
	DeviceOffline DeviceStatus = "offline"
	DeviceInGame  DeviceStatus = "in_game"
)

type GameStatus string

const (
	GameLobby    GameStatus = "lobby"
	GameRunning  GameStatus = "running"
	GameFinished GameStatus = "finished"
)

// ── Entities ──

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type Session struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Device struct {
	ID       string
	DeviceID string
	Status   DeviceStatus
	LastSeen time.Time
}

type Game struct {
	ID        string
	Code      string
	Status    GameStatus
	Settings  GameSettings
	CreatedAt time.Time
	StartedAt *time.Time
	EndedAt   *time.Time
}

// GameType defines the game mode.
// 0 = deathmatch, 1 = team deathmatch
type GameType int

const (
	GameTypeDeathmatch     GameType = 0
	GameTypeTeamDeathmatch GameType = 1
)

type GameSettings struct {
	RespawnDelay    int      `json:"respawn_delay"` // seconds
	GameDuration    int      `json:"game_duration"` // seconds, 0 = unlimited
	FriendlyFire    bool     `json:"friendly_fire"`
	MaxPlayers      int      `json:"max_players"`
	ScorePerKill    int      `json:"score_per_kill"`
	KillsPerUpgrade int      `json:"kills_per_upgrade"` // kills needed per weapon upgrade (0 = disabled)
	TypeOfGame      GameType `json:"type_of_the_game"`  // 0 = deathmatch, 1 = team deathmatch
}

func DefaultGameSettings() GameSettings {
	return GameSettings{
		RespawnDelay:    5,
		GameDuration:    300,
		FriendlyFire:    false,
		MaxPlayers:      20,
		ScorePerKill:    100,
		KillsPerUpgrade: 3,
	}
}

// Validate checks that game settings are within acceptable bounds.
func (s GameSettings) Validate() error {
	if s.MaxPlayers < 2 || s.MaxPlayers > 100 {
		return fmt.Errorf("%w: max_players must be between 2 and 100", ErrValidation)
	}
	if s.RespawnDelay < 0 || s.RespawnDelay > 300 {
		return fmt.Errorf("%w: respawn_delay must be between 0 and 300", ErrValidation)
	}
	if s.GameDuration < 0 || s.GameDuration > 7200 {
		return fmt.Errorf("%w: game_duration must be between 0 and 7200", ErrValidation)
	}
	if s.ScorePerKill < 0 || s.ScorePerKill > 10000 {
		return fmt.Errorf("%w: score_per_kill must be between 0 and 10000", ErrValidation)
	}
	if s.KillsPerUpgrade < 0 || s.KillsPerUpgrade > 50 {
		return fmt.Errorf("%w: kills_per_upgrade must be between 0 and 50", ErrValidation)
	}
	return nil
}

type Team struct {
	ID     string
	Number int
	GameID string
	Name   string
	Color  string
}

type Player struct {
	ID          string
	GameID      string
	TeamID      *string
	DeviceID    string
	Nickname    string
	Score       int
	Kills       int
	Deaths      int
	IsAlive     bool
	KillStreak  int    // consecutive kills without dying (resets on death)
	WeaponLevel int    // current weapon tier (resets on death)
	ShotsFired  int    // total shots fired this game (never resets)
	SessionCode string // unique 6-char PIN for mobile app access
}

type GameEvent struct {
	ID        string
	GameID    string
	Type      string
	Payload   map[string]any
	Timestamp time.Time
}

// ── Aggregates ──

// GameFull is the full state of a game with teams, players and events.
type GameFull struct {
	Game    Game
	Teams   []Team
	Players []Player
	Events  []GameEvent
}

// HitResult carries the outcome of a hit event.
type HitResult struct {
	Kill           bool
	AttackerID     string
	VictimID       string
	AttackerScore  int
	AttackerKills  int
	WeaponUpgraded bool // true if this kill triggered a weapon upgrade
	WeaponLevel    int  // attacker's current weapon level after this kill
	KillStreak     int  // attacker's current kill streak after this kill
}

// ReconnectInfo is returned when a reconnecting device has an active game session.
type ReconnectInfo struct {
	Player        Player
	TeamNumber    *int // nil if unassigned
	Game          Game
	RemainingTime int // seconds remaining (-1 = unlimited)
}

// PlayerSession is the public view of a player's game session for the mobile app.
type PlayerSession struct {
	Player        Player
	Game          Game
	Team          *Team // nil if unassigned
	RemainingTime int   // seconds remaining (-1 = unlimited)
}
