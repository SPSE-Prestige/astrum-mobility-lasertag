package http

import (
	"encoding/json"
	"time"
)

// Swagger DTO types — re-exported from domain for swag annotation resolution.

// ErrorResponse represents an error returned by the API.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// StatusResponse represents a simple status message.
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// TokenResponse is returned after a successful login.
type TokenResponse struct {
	Token string `json:"token" example:"abc123..."`
}

// LoginRequest holds credentials for admin login.
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"secret"`
}

// HitRequest wraps a device hit event with a game ID.
type HitRequest struct {
	GameID string         `json:"game_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Hit    HitEventDetail `json:"hit"`
}

// HitEventDetail is the hit payload from the device.
type HitEventDetail struct {
	DeviceID   string `json:"device_id" example:"vest-01"`
	TargetID   string `json:"target_id" example:"vest-02"`
	WeaponID   string `json:"weapon_id,omitempty" example:"wpn-1"`
	Damage     int    `json:"damage" example:"25"`
	IsHeadshot bool   `json:"is_headshot" example:"false"`
}

// CreateGameRequest is the body for creating a new game.
type CreateGameRequest struct {
	Name   string     `json:"name" example:"Friday Night Match"`
	Config GameConfig `json:"config"`
}

// GameConfig mirrors domain.GameConfig for Swagger.
type GameConfig struct {
	DurationSeconds int            `json:"duration_seconds" example:"300"`
	MaxPlayers      int            `json:"max_players" example:"20"`
	TeamCount       int            `json:"team_count" example:"2"`
	GameMode        string         `json:"game_mode" example:"team_deathmatch"`
	Player          PlayerConfig   `json:"player"`
	Scoring         ScoringConfig  `json:"scoring"`
	Feedback        FeedbackConfig `json:"feedback"`
}

// PlayerConfig mirrors domain.PlayerConfig.
type PlayerConfig struct {
	MaxHP               int  `json:"max_hp" example:"100"`
	Lives               int  `json:"lives" example:"3"`
	RespawnDelaySeconds int  `json:"respawn_delay_seconds" example:"5"`
	FriendlyFire        bool `json:"friendly_fire" example:"false"`
}

// ScoringConfig mirrors domain.ScoringConfig.
type ScoringConfig struct {
	PointsPerHit       int     `json:"points_per_hit" example:"10"`
	PointsPerKill      int     `json:"points_per_kill" example:"100"`
	TeamkillPenalty    int     `json:"teamkill_penalty" example:"50"`
	HeadshotMultiplier float64 `json:"headshot_multiplier" example:"2.0"`
}

// FeedbackConfig mirrors domain.FeedbackConfig.
type FeedbackConfig struct {
	SoundEnabled     bool `json:"sound_enabled" example:"true"`
	VibrationEnabled bool `json:"vibration_enabled" example:"true"`
	LEDEnabled       bool `json:"led_enabled" example:"true"`
	Intensity        int  `json:"intensity" example:"80"`
}

// GameResponse mirrors domain.Game for Swagger.
type GameResponse struct {
	ID        string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string     `json:"name" example:"Friday Night Match"`
	Status    string     `json:"status" example:"pending"`
	Config    GameConfig `json:"config"`
	CreatedAt time.Time  `json:"created_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// TeamRequest is the body for creating a team.
type TeamRequest struct {
	Name  string `json:"name" example:"Red Team"`
	Color string `json:"color" example:"#FF0000"`
}

// TeamResponse mirrors domain.Team.
type TeamResponse struct {
	ID     string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GameID string `json:"game_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name   string `json:"name" example:"Red Team"`
	Color  string `json:"color" example:"#FF0000"`
}

// JoinRequest is the body for joining a game.
type JoinRequest struct {
	Nickname string  `json:"nickname" example:"Player1"`
	DeviceID string  `json:"device_id" example:"vest-01"`
	GunID    string  `json:"gun_id" example:"gun-01"`
	UserID   *string `json:"user_id,omitempty" example:"user-uuid"`
	TeamID   *string `json:"team_id,omitempty" example:"team-uuid"`
	WeaponID *string `json:"weapon_id,omitempty" example:"wpn-1"`
}

// GamePlayerResponse mirrors domain.GamePlayer.
type GamePlayerResponse struct {
	ID             string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GameID         string  `json:"game_id"`
	UserID         *string `json:"user_id,omitempty"`
	TeamID         *string `json:"team_id,omitempty"`
	Nickname       string  `json:"nickname" example:"Player1"`
	DeviceID       string  `json:"device_id" example:"vest-01"`
	GunID          string  `json:"gun_id" example:"gun-01"`
	WeaponID       *string `json:"weapon_id,omitempty"`
	HP             int     `json:"hp" example:"100"`
	Score          int     `json:"score" example:"0"`
	Kills          int     `json:"kills" example:"0"`
	Deaths         int     `json:"deaths" example:"0"`
	IsAlive        bool    `json:"is_alive" example:"true"`
	LivesRemaining int     `json:"lives_remaining" example:"3"`
}

// LeaderboardEntryResponse mirrors domain.LeaderboardEntry.
type LeaderboardEntryResponse struct {
	PlayerID    string  `json:"player_id"`
	Nickname    string  `json:"nickname" example:"Player1"`
	TeamID      *string `json:"team_id,omitempty"`
	TeamName    string  `json:"team_name,omitempty" example:"Red Team"`
	Score       int     `json:"score" example:"350"`
	Kills       int     `json:"kills" example:"5"`
	Deaths      int     `json:"deaths" example:"2"`
	DamageDealt int     `json:"damage_dealt" example:"1250"`
}

// GameLiveStateResponse mirrors domain.GameLiveState.
type GameLiveStateResponse struct {
	GameID         string         `json:"game_id"`
	Status         string         `json:"status" example:"running"`
	TimeRemainingS int            `json:"time_remaining_s" example:"180"`
	TeamScores     map[string]int `json:"team_scores"`
}

// ControlCommandRequest mirrors domain.AdminControlCommand.
type ControlCommandRequest struct {
	Action   string                 `json:"action" example:"revive_player"`
	GameID   string                 `json:"game_id,omitempty"`
	PlayerID string                 `json:"player_id,omitempty" example:"player-uuid"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// Ensure unused imports don't fail — these are referenced
// only by Swagger annotations, not runtime code.
var (
	_ json.RawMessage
	_ time.Time
)
