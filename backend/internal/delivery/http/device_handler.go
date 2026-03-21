package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type DeviceHandler struct {
	hitUC    *usecase.HitUseCase
	gameRepo domain.GameRepository
}

func NewDeviceHandler(hitUC *usecase.HitUseCase, gameRepo domain.GameRepository) *DeviceHandler {
	return &DeviceHandler{hitUC: hitUC, gameRepo: gameRepo}
}

// HandleHit processes a hit event from a laser tag device.
// @Summary      Register a hit
// @Description  Process a hit event sent from a physical laser tag device.
// @Tags         Devices
// @Accept       json
// @Produce      json
// @Param        body body HitRequest true "Hit event"
// @Success      200 {object} StatusResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/events/hit [post]
func (h *DeviceHandler) HandleHit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameID string               `json:"game_id"`
		Hit    domain.DeviceHitEvent `json:"hit"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.GameID == "" {
		writeError(w, http.StatusBadRequest, "game_id is required")
		return
	}

	if err := h.hitUC.ProcessHit(r.Context(), req.GameID, req.Hit); err != nil {
		switch err {
		case domain.ErrNotFound:
			writeError(w, http.StatusNotFound, "game or player not found")
		case domain.ErrGameNotRunning:
			writeError(w, http.StatusConflict, "game is not running")
		case domain.ErrSamePlayer:
			writeError(w, http.StatusBadRequest, "cannot hit yourself")
		case domain.ErrNoLivesRemaining:
			writeError(w, http.StatusConflict, "player has no lives remaining")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
