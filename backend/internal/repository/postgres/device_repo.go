package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type DeviceRepo struct{ db *sql.DB }

func NewDeviceRepo(db *sql.DB) *DeviceRepo { return &DeviceRepo{db: db} }

func (r *DeviceRepo) Upsert(ctx context.Context, d *domain.Device) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`INSERT INTO devices (id, device_id, status, last_seen)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (device_id) DO UPDATE SET status = EXCLUDED.status, last_seen = EXCLUDED.last_seen`,
		d.ID, d.DeviceID, d.Status, d.LastSeen,
	)
	return err
}

func (r *DeviceRepo) GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	d := &domain.Device{}
	err := getExecutor(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, device_id, status, last_seen FROM devices WHERE device_id = $1`, deviceID,
	).Scan(&d.ID, &d.DeviceID, &d.Status, &d.LastSeen)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return d, err
}

func (r *DeviceRepo) ListAll(ctx context.Context) ([]domain.Device, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx, `SELECT id, device_id, status, last_seen FROM devices ORDER BY last_seen DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDevices(rows)
}

func (r *DeviceRepo) ListByStatus(ctx context.Context, status domain.DeviceStatus) ([]domain.Device, error) {
	rows, err := getExecutor(ctx, r.db).QueryContext(ctx,
		`SELECT id, device_id, status, last_seen FROM devices WHERE status = $1 ORDER BY last_seen DESC`, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDevices(rows)
}

func (r *DeviceRepo) UpdateStatus(ctx context.Context, deviceID string, status domain.DeviceStatus) error {
	res, err := getExecutor(ctx, r.db).ExecContext(ctx, `UPDATE devices SET status = $1 WHERE device_id = $2`, status, deviceID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DeviceRepo) UpdateLastSeen(ctx context.Context, deviceID string) error {
	// Only set status to 'online' if it was 'offline'. Preserve 'in_game' status.
	_, err := getExecutor(ctx, r.db).ExecContext(ctx,
		`UPDATE devices SET last_seen = $1, status = CASE WHEN status = $2 THEN $3 ELSE status END WHERE device_id = $4`,
		time.Now(), domain.DeviceOffline, domain.DeviceOnline, deviceID,
	)
	return err
}

func scanDevices(rows *sql.Rows) ([]domain.Device, error) {
	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		if err := rows.Scan(&d.ID, &d.DeviceID, &d.Status, &d.LastSeen); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}
