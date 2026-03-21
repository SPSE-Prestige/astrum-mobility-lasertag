package postgres

import (
	"context"
	"database/sql"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type WeaponRepo struct {
	db *sql.DB
}

func NewWeaponRepo(db *sql.DB) *WeaponRepo {
	return &WeaponRepo{db: db}
}

func (r *WeaponRepo) Create(ctx context.Context, weapon *domain.Weapon) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO weapons (id, name, damage, fire_rate_ms, ammo, reload_time_ms, fire_mode, accuracy_spread)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		weapon.ID, weapon.Name, weapon.Damage, weapon.FireRateMs,
		weapon.Ammo, weapon.ReloadTimeMs, weapon.FireMode, weapon.AccuracySpread,
	)
	return err
}

func (r *WeaponRepo) GetByID(ctx context.Context, id string) (*domain.Weapon, error) {
	var w domain.Weapon
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, damage, fire_rate_ms, ammo, reload_time_ms, fire_mode, accuracy_spread FROM weapons WHERE id = $1`, id,
	).Scan(&w.ID, &w.Name, &w.Damage, &w.FireRateMs, &w.Ammo, &w.ReloadTimeMs, &w.FireMode, &w.AccuracySpread)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &w, err
}

func (r *WeaponRepo) List(ctx context.Context) ([]*domain.Weapon, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, damage, fire_rate_ms, ammo, reload_time_ms, fire_mode, accuracy_spread FROM weapons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var weapons []*domain.Weapon
	for rows.Next() {
		var w domain.Weapon
		if err := rows.Scan(&w.ID, &w.Name, &w.Damage, &w.FireRateMs, &w.Ammo, &w.ReloadTimeMs, &w.FireMode, &w.AccuracySpread); err != nil {
			return nil, err
		}
		weapons = append(weapons, &w)
	}
	return weapons, rows.Err()
}
