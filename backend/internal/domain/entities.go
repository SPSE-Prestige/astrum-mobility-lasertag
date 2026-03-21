package domain

import "time"

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

type GameSettings struct {
	RespawnDelay int  `json:"respawn_delay"` // seconds
	GameDuration int  `json:"game_duration"` // seconds, 0 = unlimited
	FriendlyFire bool `json:"friendly_fire"`
	MaxPlayers   int  `json:"max_players"`
}

func DefaultGameSettings() GameSettings {
	return GameSettings{
		RespawnDelay: 5,
		GameDuration: 300,
		FriendlyFire: false,
		MaxPlayers:   20,
	}
}

type Team struct {
	ID     string
	GameID string
	Name   string
	Color  string
}

type Player struct {
	ID       string
	GameID   string
	TeamID   *string
	DeviceID string
	Nickname string
	Score    int
	Kills    int
	Deaths   int
	IsAlive  bool
}

type GameEvent struct {
	ID        string
	GameID    string
	Type      string
	Payload   map[string]any
	Timestamp time.Time
}
