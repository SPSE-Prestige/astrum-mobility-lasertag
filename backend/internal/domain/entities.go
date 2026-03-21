package domain

import (
	"encoding/json"
	"time"
)

// ===== Enums =====

type GameMode string

const (
	GameModeDeathmatch      GameMode = "deathmatch"
	GameModeTeamDeathmatch  GameMode = "team_deathmatch"
	GameModeCaptureTheFlag  GameMode = "capture_the_flag"
	GameModeLastManStanding GameMode = "last_man_standing"
)

type FireMode string

const (
	FireModeSingle FireMode = "single"
	FireModeBurst  FireMode = "burst"
	FireModeAuto   FireMode = "auto"
)

type GameStatus string

const (
	GameStatusPending  GameStatus = "pending"
	GameStatusRunning  GameStatus = "running"
	GameStatusPaused   GameStatus = "paused"
	GameStatusFinished GameStatus = "finished"
)

type EventType string

const (
	EventTypeHit          EventType = "hit"
	EventTypeKill         EventType = "kill"
	EventTypeRespawn      EventType = "respawn"
	EventTypeGameStart    EventType = "game_start"
	EventTypeGameEnd      EventType = "game_end"
	EventTypeGamePause    EventType = "game_pause"
	EventTypeGameResume   EventType = "game_resume"
	EventTypePlayerJoin   EventType = "player_join"
	EventTypePlayerKick   EventType = "player_kick"
	EventTypePlayerRevive EventType = "player_revive"
	EventTypeTeamChange   EventType = "team_change"
)

type AdminAction string

const (
	AdminActionPause      AdminAction = "pause_game"
	AdminActionResume     AdminAction = "resume_game"
	AdminActionEnd        AdminAction = "end_game"
	AdminActionRestart    AdminAction = "restart_game"
	AdminActionRevive     AdminAction = "revive_player"
	AdminActionKick       AdminAction = "kick_player"
	AdminActionChangeTeam AdminAction = "change_team"
)

// ===== Configuration Types =====

type PlayerConfig struct {
	MaxHP               int  `json:"max_hp"`
	Lives               int  `json:"lives"`
	RespawnDelaySeconds int  `json:"respawn_delay_seconds"`
	FriendlyFire        bool `json:"friendly_fire"`
}

type ScoringConfig struct {
	PointsPerHit       int     `json:"points_per_hit"`
	PointsPerKill      int     `json:"points_per_kill"`
	TeamkillPenalty    int     `json:"teamkill_penalty"`
	HeadshotMultiplier float64 `json:"headshot_multiplier"`
}

type FeedbackConfig struct {
	SoundEnabled     bool `json:"sound_enabled"`
	VibrationEnabled bool `json:"vibration_enabled"`
	LEDEnabled       bool `json:"led_enabled"`
	Intensity        int  `json:"intensity"`
}

type GameConfig struct {
	DurationSeconds int            `json:"duration_seconds"`
	MaxPlayers      int            `json:"max_players"`
	TeamCount       int            `json:"team_count"`
	GameMode        GameMode       `json:"game_mode"`
	Player          PlayerConfig   `json:"player"`
	Scoring         ScoringConfig  `json:"scoring"`
	Feedback        FeedbackConfig `json:"feedback"`
}

// ===== Entities =====

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Game struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    GameStatus `json:"status"`
	Config    GameConfig `json:"config"`
	CreatedAt time.Time  `json:"created_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

func (g *Game) ConfigJSON() ([]byte, error) {
	return json.Marshal(g.Config)
}

type Team struct {
	ID     string `json:"id"`
	GameID string `json:"game_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type GamePlayer struct {
	ID             string  `json:"id"`
	GameID         string  `json:"game_id"`
	UserID         *string `json:"user_id,omitempty"`
	TeamID         *string `json:"team_id,omitempty"`
	Nickname       string  `json:"nickname"`
	DeviceID       string  `json:"device_id"`
	GunID          string  `json:"gun_id"`
	WeaponID       *string `json:"weapon_id,omitempty"`
	HP             int     `json:"hp"`
	Score          int     `json:"score"`
	Kills          int     `json:"kills"`
	Deaths         int     `json:"deaths"`
	IsAlive        bool    `json:"is_alive"`
	LivesRemaining int     `json:"lives_remaining"`
}

type Weapon struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Damage         int      `json:"damage"`
	FireRateMs     int      `json:"fire_rate_ms"`
	Ammo           int      `json:"ammo"`
	ReloadTimeMs   int      `json:"reload_time_ms"`
	FireMode       FireMode `json:"fire_mode"`
	AccuracySpread float64  `json:"accuracy_spread"`
}

type GameEvent struct {
	ID        string          `json:"id"`
	GameID    string          `json:"game_id"`
	Type      EventType       `json:"type"`
	PlayerID  string          `json:"player_id"`
	TargetID  string          `json:"target_id"`
	WeaponID  string          `json:"weapon_id"`
	Damage    int             `json:"damage"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type AdminSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ===== Device Events =====

type DeviceHitEvent struct {
	DeviceID   string `json:"device_id"`
	TargetID   string `json:"target_id"`
	WeaponID   string `json:"weapon_id,omitempty"`
	Damage     int    `json:"damage"`
	IsHeadshot bool   `json:"is_headshot"`
}

// ===== Live State (cached in Redis) =====

type PlayerLiveState struct {
	PlayerID       string `json:"player_id"`
	GameID         string `json:"game_id"`
	HP             int    `json:"hp"`
	Score          int    `json:"score"`
	Kills          int    `json:"kills"`
	Deaths         int    `json:"deaths"`
	IsAlive        bool   `json:"is_alive"`
	LivesRemaining int    `json:"lives_remaining"`
}

type GameLiveState struct {
	GameID         string         `json:"game_id"`
	Status         GameStatus     `json:"status"`
	TimeRemainingS int            `json:"time_remaining_s"`
	TeamScores     map[string]int `json:"team_scores"`
}

// ===== Leaderboard =====

type LeaderboardEntry struct {
	PlayerID    string  `json:"player_id"`
	Nickname    string  `json:"nickname"`
	TeamID      *string `json:"team_id,omitempty"`
	TeamName    string  `json:"team_name,omitempty"`
	Score       int     `json:"score"`
	Kills       int     `json:"kills"`
	Deaths      int     `json:"deaths"`
	DamageDealt int     `json:"damage_dealt"`
}

// ===== Admin Control =====

type AdminControlCommand struct {
	Action   AdminAction            `json:"action"`
	GameID   string                 `json:"game_id"`
	PlayerID string                 `json:"player_id,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}
