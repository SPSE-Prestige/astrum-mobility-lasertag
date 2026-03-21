package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrSessionExpired   = errors.New("session expired")
	ErrInvalidGameState = errors.New("invalid game state")
	ErrPlayerDead       = errors.New("player is dead")
	ErrFriendlyFire     = errors.New("friendly fire not allowed")
	ErrSelfHit          = errors.New("cannot hit yourself")
	ErrGameFull         = errors.New("game is full")
	ErrDeviceInGame     = errors.New("device already in a game")
)
