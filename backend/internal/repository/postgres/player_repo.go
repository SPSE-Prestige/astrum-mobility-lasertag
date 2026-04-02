package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type PlayerRepo struct{ db *sql.DB }

func NewPlayerRepo(db *sql.DB) *PlayerRepo { return &PlayerRepo{db: db} }

func (r *PlayerRepo) Create(ctx context.Context, p *domain.Player) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`INSERT INTO players (id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		p.ID, p.GameID, nilIfEmpty(p.TeamID), p.DeviceID, p.Nickname,
		p.Score, p.Kills, p.Deaths, p.IsAlive, p.KillStreak, p.WeaponLevel, p.ShotsFired, nilIfEmptyStr(p.SessionCode),
	)
	return err
}

func (r *PlayerRepo) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	return r.scanPlayer(getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code FROM players WHERE id = $1`, id,
	))
}

func (r *PlayerRepo) GetByGameAndDevice(ctx context.Context, gameID, deviceID string) (*domain.Player, error) {
	return r.scanPlayer(getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code
		 FROM players WHERE game_id = $1 AND device_id = $2`, gameID, deviceID,
	))
}

// FindActivePlayerByDevice finds a player with the given device in a running game.
func (r *PlayerRepo) FindActivePlayerByDevice(ctx context.Context, deviceID string) (*domain.Player, *domain.Game, error) {
	row := getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT p.id, p.game_id, p.team_id, p.device_id, p.nickname, p.score, p.kills, p.deaths, p.is_alive, p.kill_streak, p.weapon_level, p.shots_fired, p.session_code,
		        g.id, g.code, g.status, g.settings, g.created_at, g.started_at, g.ended_at
		 FROM players p
		 JOIN games g ON p.game_id = g.id
		 WHERE p.device_id = $1 AND g.status = $2
		 LIMIT 1`, deviceID, domain.GameRunning,
	)

	p := &domain.Player{}
	g := &domain.Game{}
	var teamID, sessionCode sql.NullString
	var startedAt, endedAt sql.NullTime
	var settingsJSON []byte

	err := row.Scan(
		&p.ID, &p.GameID, &teamID, &p.DeviceID, &p.Nickname, &p.Score, &p.Kills, &p.Deaths, &p.IsAlive, &p.KillStreak, &p.WeaponLevel, &p.ShotsFired, &sessionCode,
		&g.ID, &g.Code, &g.Status, &settingsJSON, &g.CreatedAt, &startedAt, &endedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil // not in any active game
	}
	if err != nil {
		return nil, nil, fmt.Errorf("find active player by device %s: %w", deviceID, err)
	}

	_ = json.Unmarshal(settingsJSON, &g.Settings)
	if teamID.Valid {
		p.TeamID = &teamID.String
	}
	if sessionCode.Valid {
		p.SessionCode = sessionCode.String
	}
	if startedAt.Valid {
		g.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		g.EndedAt = &endedAt.Time
	}
	return p, g, nil
}

func (r *PlayerRepo) ListByGame(ctx context.Context, gameID string) ([]domain.Player, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code FROM players WHERE game_id = $1 ORDER BY score DESC`, gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPlayers(rows)
}

func (r *PlayerRepo) ListByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code FROM players WHERE team_id = $1 ORDER BY score DESC`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPlayers(rows)
}

func (r *PlayerRepo) Update(ctx context.Context, p *domain.Player) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE players SET team_id=$1, nickname=$2, score=$3, kills=$4, deaths=$5, is_alive=$6, kill_streak=$7, weapon_level=$8, shots_fired=$9, session_code=$10 WHERE id=$11`,
		nilIfEmpty(p.TeamID), p.Nickname, p.Score, p.Kills, p.Deaths, p.IsAlive, p.KillStreak, p.WeaponLevel, p.ShotsFired, nilIfEmptyStr(p.SessionCode), p.ID,
	)
	return err
}

func (r *PlayerRepo) Delete(ctx context.Context, id string) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM players WHERE id = $1`, id)
	return err
}

// KillPlayer atomically sets is_alive=false, increments deaths, and resets kill_streak + weapon_level.
// Returns false if the player was already dead (prevents duplicate kills).
func (r *PlayerRepo) KillPlayer(ctx context.Context, playerID string) (bool, error) {
	res, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE players SET is_alive = false, deaths = deaths + 1, kill_streak = 0, weapon_level = 0
		 WHERE id = $1 AND is_alive = true`,
		playerID,
	)
	if err != nil {
		return false, fmt.Errorf("kill player %s: %w", playerID, err)
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// AddKillScore atomically increments kills, score, and kill_streak.
// If killsPerUpgrade > 0 and the new streak hits the threshold, weapon_level is also incremented.
// Returns the updated streak state.
func (r *PlayerRepo) AddKillScore(ctx context.Context, playerID string, score, killsPerUpgrade int) (*domain.KillScoreResult, error) {
	row := getExecutor(ctx, r.db).QueryRowContext(ctx,
		`UPDATE players SET
			kills = kills + 1,
			score = score + $1,
			kill_streak = kill_streak + 1,
			weapon_level = CASE
				WHEN $2 > 0 AND (kill_streak + 1) % $2 = 0 THEN weapon_level + 1
				ELSE weapon_level
			END
		 WHERE id = $3
		 RETURNING kill_streak, weapon_level`,
		score, killsPerUpgrade, playerID,
	)
	var result domain.KillScoreResult
	if err := row.Scan(&result.KillStreak, &result.WeaponLevel); err != nil {
		return nil, fmt.Errorf("add kill score %s: %w", playerID, err)
	}
	return &result, nil
}

// Respawn atomically sets is_alive=true.
func (r *PlayerRepo) Respawn(ctx context.Context, playerID string) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE players SET is_alive = true WHERE id = $1`,
		playerID,
	)
	return err
}

// IncrementShotsFired atomically increments the shots_fired counter.
func (r *PlayerRepo) IncrementShotsFired(ctx context.Context, playerID string) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE players SET shots_fired = shots_fired + 1 WHERE id = $1`,
		playerID,
	)
	return err
}

func (r *PlayerRepo) scanPlayer(row *sql.Row) (*domain.Player, error) {
	p := &domain.Player{}
	var teamID, sessionCode sql.NullString
	err := row.Scan(&p.ID, &p.GameID, &teamID, &p.DeviceID, &p.Nickname,
		&p.Score, &p.Kills, &p.Deaths, &p.IsAlive, &p.KillStreak, &p.WeaponLevel, &p.ShotsFired, &sessionCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if teamID.Valid {
		p.TeamID = &teamID.String
	}
	if sessionCode.Valid {
		p.SessionCode = sessionCode.String
	}
	return p, nil
}

func (r *PlayerRepo) scanPlayers(rows *sql.Rows) ([]domain.Player, error) {
	var players []domain.Player
	for rows.Next() {
		var p domain.Player
		var teamID, sessionCode sql.NullString
		if err := rows.Scan(&p.ID, &p.GameID, &teamID, &p.DeviceID, &p.Nickname,
			&p.Score, &p.Kills, &p.Deaths, &p.IsAlive, &p.KillStreak, &p.WeaponLevel, &p.ShotsFired, &sessionCode); err != nil {
			return nil, err
		}
		if teamID.Valid {
			p.TeamID = &teamID.String
		}
		if sessionCode.Valid {
			p.SessionCode = sessionCode.String
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

// GetBySessionCode finds a player by their unique session PIN code.
func (r *PlayerRepo) GetBySessionCode(ctx context.Context, code string) (*domain.Player, error) {
	return r.scanPlayer(getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, game_id, team_id, device_id, nickname, score, kills, deaths, is_alive, kill_streak, weapon_level, shots_fired, session_code
		 FROM players WHERE session_code = $1`, code,
	))
}

func nilIfEmpty(s *string) any {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func nilIfEmptyStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
