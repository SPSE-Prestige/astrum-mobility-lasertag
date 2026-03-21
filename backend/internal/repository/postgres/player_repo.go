package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type GamePlayerRepo struct {
	db *sql.DB
}

func NewGamePlayerRepo(db *sql.DB) *GamePlayerRepo {
	return &GamePlayerRepo{db: db}
}

func (r *GamePlayerRepo) Create(ctx context.Context, player *domain.GamePlayer) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO game_players (id, game_id, user_id, team_id, nickname, device_id, gun_id, weapon_id, hp, score, kills, deaths, is_alive, lives_remaining)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		player.ID, player.GameID, player.UserID, player.TeamID, player.Nickname,
		player.DeviceID, player.GunID, player.WeaponID, player.HP, player.Score,
		player.Kills, player.Deaths, player.IsAlive, player.LivesRemaining,
	)
	return err
}

func (r *GamePlayerRepo) GetByID(ctx context.Context, id string) (*domain.GamePlayer, error) {
	return r.scanOne(r.db.QueryRowContext(ctx,
		`SELECT id, game_id, user_id, team_id, nickname, device_id, gun_id, weapon_id, hp, score, kills, deaths, is_alive, lives_remaining
		 FROM game_players WHERE id = $1`, id))
}

func (r *GamePlayerRepo) GetByDeviceID(ctx context.Context, gameID, deviceID string) (*domain.GamePlayer, error) {
	return r.scanOne(r.db.QueryRowContext(ctx,
		`SELECT id, game_id, user_id, team_id, nickname, device_id, gun_id, weapon_id, hp, score, kills, deaths, is_alive, lives_remaining
		 FROM game_players WHERE game_id = $1 AND device_id = $2`, gameID, deviceID))
}

func (r *GamePlayerRepo) GetByGunID(ctx context.Context, gameID, gunID string) (*domain.GamePlayer, error) {
	return r.scanOne(r.db.QueryRowContext(ctx,
		`SELECT id, game_id, user_id, team_id, nickname, device_id, gun_id, weapon_id, hp, score, kills, deaths, is_alive, lives_remaining
		 FROM game_players WHERE game_id = $1 AND gun_id = $2`, gameID, gunID))
}

func (r *GamePlayerRepo) GetByGameID(ctx context.Context, gameID string) ([]*domain.GamePlayer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, game_id, user_id, team_id, nickname, device_id, gun_id, weapon_id, hp, score, kills, deaths, is_alive, lives_remaining
		 FROM game_players WHERE game_id = $1`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []*domain.GamePlayer
	for rows.Next() {
		p, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func (r *GamePlayerRepo) Update(ctx context.Context, player *domain.GamePlayer) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE game_players SET user_id=$1, team_id=$2, nickname=$3, device_id=$4, gun_id=$5, weapon_id=$6,
		 hp=$7, score=$8, kills=$9, deaths=$10, is_alive=$11, lives_remaining=$12 WHERE id=$13`,
		player.UserID, player.TeamID, player.Nickname, player.DeviceID, player.GunID, player.WeaponID,
		player.HP, player.Score, player.Kills, player.Deaths, player.IsAlive, player.LivesRemaining, player.ID,
	)
	return err
}

func (r *GamePlayerRepo) CountByGameID(ctx context.Context, gameID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM game_players WHERE game_id = $1`, gameID).Scan(&count)
	return count, err
}

func (r *GamePlayerRepo) scanOne(row *sql.Row) (*domain.GamePlayer, error) {
	var p domain.GamePlayer
	err := row.Scan(
		&p.ID, &p.GameID, &p.UserID, &p.TeamID, &p.Nickname,
		&p.DeviceID, &p.GunID, &p.WeaponID, &p.HP, &p.Score,
		&p.Kills, &p.Deaths, &p.IsAlive, &p.LivesRemaining,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &p, err
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *GamePlayerRepo) scanRow(row rowScanner) (*domain.GamePlayer, error) {
	var p domain.GamePlayer
	err := row.Scan(
		&p.ID, &p.GameID, &p.UserID, &p.TeamID, &p.Nickname,
		&p.DeviceID, &p.GunID, &p.WeaponID, &p.HP, &p.Score,
		&p.Kills, &p.Deaths, &p.IsAlive, &p.LivesRemaining,
	)
	return &p, err
}
