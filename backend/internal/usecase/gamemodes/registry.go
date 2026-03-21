package gamemodes

import (
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

// Registry stores all available game mode handlers.
type Registry struct {
	modes map[domain.GameMode]domain.GameModeHandler
}

func NewRegistry() *Registry {
	r := &Registry{modes: make(map[domain.GameMode]domain.GameModeHandler)}
	r.Register(&Deathmatch{})
	r.Register(&TeamDeathmatch{})
	r.Register(&LastManStanding{})
	return r
}

func (r *Registry) Register(handler domain.GameModeHandler) {
	r.modes[handler.Mode()] = handler
}

func (r *Registry) Get(mode domain.GameMode) (domain.GameModeHandler, error) {
	h, ok := r.modes[mode]
	if !ok {
		return nil, domain.ErrInvalidGameMode
	}
	return h, nil
}
