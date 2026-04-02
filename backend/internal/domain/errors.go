package domain

import "errors"

// ── Sentinel errors ──

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
	ErrValidation       = errors.New("validation error")
	ErrInternal         = errors.New("internal error")
)

// ── API Error with code ──

// AppError represents a structured application error with a machine-readable code.
type AppError struct {
	Code    string // machine-readable e.g. "GAME_NOT_FOUND"
	Message string // human-readable
	Err     error  // underlying error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
