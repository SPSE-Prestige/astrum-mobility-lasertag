package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

const maxBodySize = 1 << 20 // 1 MB

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, ErrorResponse{Error: ErrorDetail{Code: code, Message: msg}})
}

func readJSON(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func fmtTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func fmtTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func toDeviceResponse(d domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:       d.ID,
		DeviceID: d.DeviceID,
		Status:   string(d.Status),
		LastSeen: fmtTime(d.LastSeen),
	}
}

func toGameResponse(g domain.Game) GameResponse {
	return GameResponse{
		ID:     g.ID,
		Code:   g.Code,
		Status: string(g.Status),
		Settings: GameSettingsDTO{
			RespawnDelay:    g.Settings.RespawnDelay,
			GameDuration:    g.Settings.GameDuration,
			FriendlyFire:    g.Settings.FriendlyFire,
			MaxPlayers:      g.Settings.MaxPlayers,
			ScorePerKill:    g.Settings.ScorePerKill,
			KillsPerUpgrade: g.Settings.KillsPerUpgrade,
		},
		CreatedAt: fmtTime(g.CreatedAt),
		StartedAt: fmtTimePtr(g.StartedAt),
		EndedAt:   fmtTimePtr(g.EndedAt),
	}
}

func toTeamResponse(t domain.Team) TeamResponse {
	return TeamResponse{
		ID:     t.ID,
		GameID: t.GameID,
		Name:   t.Name,
		Color:  t.Color,
	}
}

func toPlayerResponse(p domain.Player) PlayerResponse {
	return PlayerResponse{
		ID:          p.ID,
		GameID:      p.GameID,
		TeamID:      p.TeamID,
		DeviceID:    p.DeviceID,
		Nickname:    p.Nickname,
		Score:       p.Score,
		Kills:       p.Kills,
		Deaths:      p.Deaths,
		IsAlive:     p.IsAlive,
		KillStreak:  p.KillStreak,
		WeaponLevel: p.WeaponLevel,
		ShotsFired:  p.ShotsFired,
	}
}

func toEventResponse(e domain.GameEvent) EventResponse {
	return EventResponse{
		ID:        e.ID,
		GameID:    e.GameID,
		Type:      e.Type,
		Payload:   e.Payload,
		Timestamp: fmtTime(e.Timestamp),
	}
}

func mapSlice[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}
