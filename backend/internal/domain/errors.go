package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrGameNotPending     = errors.New("game is not in pending state")
	ErrGameNotRunning     = errors.New("game is not running")
	ErrGameNotPaused      = errors.New("game is not paused")
	ErrGameAlreadyRunning = errors.New("game is already running")
	ErrPlayerDead         = errors.New("player is dead")
	ErrSameTeam           = errors.New("cannot hit teammate")
	ErrSamePlayer         = errors.New("cannot hit yourself")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrSessionExpired     = errors.New("session expired")
	ErrDuplicateDevice    = errors.New("device already assigned in this game")
	ErrGameFull           = errors.New("game is full")
	ErrNoLivesRemaining   = errors.New("no lives remaining")
	ErrInvalidGameMode    = errors.New("invalid game mode")
	ErrInvalidAction      = errors.New("invalid admin action")
	ErrPlayerNotInGame    = errors.New("player not in this game")
	ErrWeaponCooldown     = errors.New("weapon on cooldown")
	ErrNoAmmo             = errors.New("no ammo remaining")
)
