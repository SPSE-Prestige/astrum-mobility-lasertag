package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type PlayerHandler struct {
	playerUC *usecase.PlayerUseCase
}

func NewPlayerHandler(playerUC *usecase.PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{playerUC: playerUC}
}

// Join adds a player to a game.
// @Summary      Join a game
// @Description  Register a player (with device and gun) into an existing game.
// @Tags         Players
// @Accept       json
// @Produce      json
// @Param        id   path string true "Game ID"
// @Param        body body JoinRequest true "Player data"
// @Success      201 {object} GamePlayerResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/join [post]
func (h *PlayerHandler) Join(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	var req struct {
		Nickname string  `json:"nickname"`
		DeviceID string  `json:"device_id"`
		GunID    string  `json:"gun_id"`
		UserID   *string `json:"user_id,omitempty"`
		TeamID   *string `json:"team_id,omitempty"`
		WeaponID *string `json:"weapon_id,omitempty"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Nickname == "" || req.DeviceID == "" || req.GunID == "" {
		writeError(w, http.StatusBadRequest, "nickname, device_id, and gun_id are required")
		return
	}

	player, err := h.playerUC.JoinGame(r.Context(), gameID, req.Nickname, req.DeviceID, req.GunID, req.UserID, req.TeamID, req.WeaponID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, player)
}

// Leaderboard returns the leaderboard for a game.
// @Summary      Get leaderboard
// @Description  Retrieve sorted leaderboard entries for a game.
// @Tags         Players
// @Produce      json
// @Param        id path string true "Game ID"
// @Success      200 {array} LeaderboardEntryResponse
// @Failure      404 {object} ErrorResponse
// @Router       /api/games/{id}/leaderboard [get]
func (h *PlayerHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	entries, err := h.playerUC.GetLeaderboard(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// ListPlayers returns all players in a game.
// @Summary      List players
// @Description  Get all players currently registered in a game.
// @Tags         Players
// @Produce      json
// @Param        id path string true "Game ID"
// @Success      200 {array} GamePlayerResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/games/{id}/players [get]
func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	players, err := h.playerUC.GetPlayers(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, players)
}
