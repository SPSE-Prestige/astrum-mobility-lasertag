package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type GameHandler struct {
	gameUC  *usecase.GameUseCase
	adminUC *usecase.AdminUseCase
}

func NewGameHandler(gameUC *usecase.GameUseCase, adminUC *usecase.AdminUseCase) *GameHandler {
	return &GameHandler{gameUC: gameUC, adminUC: adminUC}
}

// Create creates a new game.
// @Summary      Create a game
// @Description  Create a new laser tag game with the given configuration.
// @Tags         Games
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreateGameRequest true "Game definition"
// @Success      201 {object} GameResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games [post]
func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Config domain.GameConfig `json:"config"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	game, err := h.gameUC.CreateGame(r.Context(), req.Name, req.Config)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, game)
}

// Get returns a single game by ID.
// @Summary      Get a game
// @Description  Retrieve a game by its ID.
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Game ID"
// @Success      200 {object} GameResponse
// @Failure      404 {object} ErrorResponse
// @Router       /api/games/{id} [get]
func (h *GameHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	game, err := h.gameUC.GetGame(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	writeJSON(w, http.StatusOK, game)
}

// List returns all games.
// @Summary      List games
// @Description  Get a list of all games.
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} GameResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/games [get]
func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameUC.ListGames(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, games)
}

// Start starts a game.
// @Summary      Start a game
// @Description  Transition a game from pending to running.
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Game ID"
// @Success      200 {object} StatusResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/start [post]
func (h *GameHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.gameUC.StartGame(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

// Pause pauses a running game.
// @Summary      Pause a game
// @Description  Pause a currently running game.
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Game ID"
// @Success      200 {object} StatusResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/pause [post]
func (h *GameHandler) Pause(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.gameUC.PauseGame(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
}

// End finishes a game.
// @Summary      End a game
// @Description  End a running or paused game and persist final state.
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Game ID"
// @Success      200 {object} StatusResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/end [post]
func (h *GameHandler) End(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.gameUC.EndGame(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

// Control executes an admin control command on a game.
// @Summary      Admin control
// @Description  Execute an admin action (revive, kick, change_team, restart, etc.).
// @Tags         Games
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                    true "Game ID"
// @Param        body body ControlCommandRequest true "Control command"
// @Success      200 {object} StatusResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/control [post]
func (h *GameHandler) Control(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd domain.AdminControlCommand
	if err := readJSON(r, &cmd); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cmd.GameID = id
	if err := h.adminUC.ExecuteControl(r.Context(), cmd); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// CreateTeam adds a team to a game.
// @Summary      Create a team
// @Description  Add a new team to an existing game.
// @Tags         Games
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                            true "Game ID"
// @Param        body body TeamRequest true "Team definition"
// @Success      201 {object} TeamResponse
// @Failure      400 {object} ErrorResponse
// @Router       /api/games/{id}/teams [post]
func (h *GameHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	team, err := h.gameUC.CreateTeam(r.Context(), id, req.Name, req.Color)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, team)
}

// GetState returns the live state of a game.
// @Summary      Get game state
// @Description  Retrieve the live state of a game from cache (players, scores, timer).
// @Tags         Games
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Game ID"
// @Success      200 {object} GameLiveStateResponse
// @Failure      404 {object} ErrorResponse
// @Router       /api/games/{id}/state [get]
func (h *GameHandler) GetState(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	state, err := h.gameUC.GetGameState(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "game state not found")
		return
	}
	writeJSON(w, http.StatusOK, state)
}
