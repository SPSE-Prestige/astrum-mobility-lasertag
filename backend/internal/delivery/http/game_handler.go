package http

import (
	"errors"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/infrastructure/mqtt"
	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type GameHandler struct {
	gameUC *usecase.GameUseCase
	mqtt   *mqtt.Client
}

func NewGameHandler(gameUC *usecase.GameUseCase, mqttClient *mqtt.Client) *GameHandler {
	return &GameHandler{gameUC: gameUC, mqtt: mqttClient}
}

// Create godoc
//
//	@Summary	Create a new game
//	@Tags		games
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		CreateGameRequest	true	"Game settings (optional)"
//	@Success	201	{object}	GameResponse
//	@Failure	400	{object}	ErrorResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games [post]
func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var settings *domain.GameSettings
	if req.Settings != nil {
		s := domain.GameSettings{
			RespawnDelay: req.Settings.RespawnDelay,
			GameDuration: req.Settings.GameDuration,
			FriendlyFire: req.Settings.FriendlyFire,
			MaxPlayers:   req.Settings.MaxPlayers,
		}
		settings = &s
	}

	game, err := h.gameUC.CreateGame(r.Context(), settings)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toGameResponse(*game))
}

// Get godoc
//
//	@Summary	Get game by ID
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{object}	GameResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Router		/games/{id} [get]
func (h *GameHandler) Get(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	game, err := h.gameUC.GetGame(r.Context(), gameID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toGameResponse(*game))
}

// GetFull godoc
//
//	@Summary	Get full game state (game + teams + players + events)
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{object}	GameFullResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Router		/games/{id}/full [get]
func (h *GameHandler) GetFull(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	full, err := h.gameUC.GetGameFull(r.Context(), gameID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, GameFullResponse{
		Game:    toGameResponse(full.Game),
		Teams:   mapSlice(full.Teams, toTeamResponse),
		Players: mapSlice(full.Players, toPlayerResponse),
		Events:  mapSlice(full.Events, toEventResponse),
	})
}

// List godoc
//
//	@Summary	List all games
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}		GameResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games [get]
func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameUC.ListGames(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(games, toGameResponse))
}

// Start godoc
//
//	@Summary	Start a game (lobby → running)
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{object}	GameResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	409	{object}	ErrorResponse
//	@Router		/games/{id}/start [post]
func (h *GameHandler) Start(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	game, err := h.gameUC.StartGame(r.Context(), gameID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrInvalidGameState) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}

	// Notify devices via MQTT
	players, _ := h.gameUC.ListPlayers(r.Context(), gameID)
	deviceIDs := make([]string, len(players))
	for i, p := range players {
		deviceIDs[i] = p.DeviceID
	}
	h.mqtt.PublishGameStart(deviceIDs, gameID)

	writeJSON(w, http.StatusOK, toGameResponse(*game))
}

// End godoc
//
//	@Summary	End a game (running → finished)
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{object}	GameResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	409	{object}	ErrorResponse
//	@Router		/games/{id}/end [post]
func (h *GameHandler) End(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	game, err := h.gameUC.EndGame(r.Context(), gameID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrInvalidGameState) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}

	// Notify devices via MQTT
	players, _ := h.gameUC.ListPlayers(r.Context(), gameID)
	deviceIDs := make([]string, len(players))
	for i, p := range players {
		deviceIDs[i] = p.DeviceID
	}
	h.mqtt.PublishGameEnd(deviceIDs)

	writeJSON(w, http.StatusOK, toGameResponse(*game))
}

// UpdateSettings godoc
//
//	@Summary	Update game settings (lobby only)
//	@Tags		games
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string				true	"Game ID"
//	@Param		body	body		UpdateSettingsRequest	true	"New settings"
//	@Success	200		{object}	GameResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/games/{id}/settings [patch]
func (h *GameHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	var req UpdateSettingsRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	settings := domain.GameSettings{
		RespawnDelay: req.Settings.RespawnDelay,
		GameDuration: req.Settings.GameDuration,
		FriendlyFire: req.Settings.FriendlyFire,
		MaxPlayers:   req.Settings.MaxPlayers,
	}

	game, err := h.gameUC.UpdateSettings(r.Context(), gameID, settings)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrInvalidGameState) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toGameResponse(*game))
}

// ── Team endpoints ──

// AddTeam godoc
//
//	@Summary	Add a team to a game
//	@Tags		teams
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string			true	"Game ID"
//	@Param		body	body		AddTeamRequest	true	"Team data"
//	@Success	201		{object}	TeamResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/games/{id}/teams [post]
func (h *GameHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	var req AddTeamRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Color == "" {
		writeError(w, http.StatusBadRequest, "name and color required")
		return
	}

	team, err := h.gameUC.AddTeam(r.Context(), gameID, req.Name, req.Color)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toTeamResponse(*team))
}

// ListTeams godoc
//
//	@Summary	List teams in a game
//	@Tags		teams
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{array}		TeamResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/teams [get]
func (h *GameHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	teams, err := h.gameUC.ListTeams(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(teams, toTeamResponse))
}

// RemoveTeam godoc
//
//	@Summary	Remove a team
//	@Tags		teams
//	@Security	BearerAuth
//	@Param		id		path	string	true	"Game ID"
//	@Param		teamId	path	string	true	"Team ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/teams/{teamId} [delete]
func (h *GameHandler) RemoveTeam(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if err := h.gameUC.RemoveTeam(r.Context(), teamID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Player endpoints ──

// AddPlayer godoc
//
//	@Summary	Add a player (device) to a game
//	@Tags		players
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string				true	"Game ID"
//	@Param		body	body		AddPlayerRequest	true	"Player data"
//	@Success	201		{object}	PlayerResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse	"Game full or device already in game"
//	@Router		/games/{id}/players [post]
func (h *GameHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	var req AddPlayerRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DeviceID == "" || req.Nickname == "" {
		writeError(w, http.StatusBadRequest, "device_id and nickname required")
		return
	}

	player, err := h.gameUC.AddPlayer(r.Context(), gameID, req.DeviceID, req.Nickname, req.TeamID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrGameFull) {
			status = http.StatusConflict
		} else if errors.Is(err, domain.ErrDeviceInGame) {
			status = http.StatusConflict
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toPlayerResponse(*player))
}

// ListPlayers godoc
//
//	@Summary	List players in a game
//	@Tags		players
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{array}		PlayerResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/players [get]
func (h *GameHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	players, err := h.gameUC.ListPlayers(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(players, toPlayerResponse))
}

// RemovePlayer godoc
//
//	@Summary	Remove a player from a game
//	@Tags		players
//	@Security	BearerAuth
//	@Param		id			path	string	true	"Game ID"
//	@Param		playerId	path	string	true	"Player ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/players/{playerId} [delete]
func (h *GameHandler) RemovePlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	if err := h.gameUC.RemovePlayer(r.Context(), playerID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdatePlayerTeam godoc
//
//	@Summary	Change a player's team assignment
//	@Tags		players
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id			path	string						true	"Game ID"
//	@Param		playerId	path	string						true	"Player ID"
//	@Param		body		body	UpdatePlayerTeamRequest		true	"New team"
//	@Success	204
//	@Failure	400	{object}	ErrorResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	409	{object}	ErrorResponse
//	@Router		/games/{id}/players/{playerId}/team [patch]
func (h *GameHandler) UpdatePlayerTeam(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	var req UpdatePlayerTeamRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.gameUC.UpdatePlayerTeam(r.Context(), playerID, req.TeamID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrInvalidGameState) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Leaderboard & Events ──

// Leaderboard godoc
//
//	@Summary	Get game leaderboard (players sorted by score)
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{array}		PlayerResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/leaderboard [get]
func (h *GameHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	players, err := h.gameUC.GetLeaderboard(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(players, toPlayerResponse))
}

// Events godoc
//
//	@Summary	Get game events (kills, respawns, etc.)
//	@Tags		games
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Game ID"
//	@Success	200	{array}		EventResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/games/{id}/events [get]
func (h *GameHandler) Events(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	events, err := h.gameUC.ListEvents(r.Context(), gameID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(events, toEventResponse))
}
