package http

import "time"

// ── Request DTOs ──

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateGameRequest struct {
	Settings *GameSettingsDTO `json:"settings,omitempty"`
}

type GameSettingsDTO struct {
	RespawnDelay int  `json:"respawn_delay"`
	GameDuration int  `json:"game_duration"`
	FriendlyFire bool `json:"friendly_fire"`
	MaxPlayers   int  `json:"max_players"`
}

type AddTeamRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type AddPlayerRequest struct {
	DeviceID string  `json:"device_id"`
	Nickname string  `json:"nickname"`
	TeamID   *string `json:"team_id,omitempty"`
}

type UpdatePlayerTeamRequest struct {
	TeamID *string `json:"team_id"`
}

type UpdateSettingsRequest struct {
	Settings GameSettingsDTO `json:"settings"`
}

// ── Response DTOs ──

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DeviceResponse struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen"`
}

type GameResponse struct {
	ID        string          `json:"id"`
	Code      string          `json:"code"`
	Status    string          `json:"status"`
	Settings  GameSettingsDTO `json:"settings"`
	CreatedAt string          `json:"created_at"`
	StartedAt *string         `json:"started_at,omitempty"`
	EndedAt   *string         `json:"ended_at,omitempty"`
}

type TeamResponse struct {
	ID     string `json:"id"`
	GameID string `json:"game_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type PlayerResponse struct {
	ID       string  `json:"id"`
	GameID   string  `json:"game_id"`
	TeamID   *string `json:"team_id,omitempty"`
	DeviceID string  `json:"device_id"`
	Nickname string  `json:"nickname"`
	Score    int     `json:"score"`
	Kills    int     `json:"kills"`
	Deaths   int     `json:"deaths"`
	IsAlive  bool    `json:"is_alive"`
}

type GameFullResponse struct {
	Game    GameResponse     `json:"game"`
	Teams   []TeamResponse   `json:"teams"`
	Players []PlayerResponse `json:"players"`
	Events  []EventResponse  `json:"events"`
}

type EventResponse struct {
	ID        string         `json:"id"`
	GameID    string         `json:"game_id"`
	Type      string         `json:"type"`
	Payload   map[string]any `json:"payload"`
	Timestamp string         `json:"timestamp"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
