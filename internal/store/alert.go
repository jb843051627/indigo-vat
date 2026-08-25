package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutAlert(ctx context.Context, a model.Alert) error {
	_, err := d.Exec(ctx, `INSERT INTO alerts(id,cycle_id,level,code,message,state,created_at,acknowledged_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.CycleID, a.Level, a.Code, a.Message, a.State, Text(a.CreatedAt), Text(a.AcknowledgedAt))
	return err
}
func (d *DB) ListAlerts(ctx context.Context, cycleID string) ([]model.Alert, error) {
	rows, err := d.Query(ctx, `SELECT id,cycle_id,level,code,message,state,created_at,acknowledged_at FROM alerts WHERE cycle_id=? ORDER BY created_at`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Alert{}
	for rows.Next() {
		var a model.Alert
		var c, d string
		if err := rows.Scan(&a.ID, &a.CycleID, &a.Level, &a.Code, &a.Message, &a.State, &c, &d); err != nil {
			return nil, err
		}
		a.CreatedAt = Time(c)
		a.AcknowledgedAt = Time(d)
		out = append(out, a)
	}
	return out, rows.Err()
}
func (d *DB) OpenAlertCount(ctx context.Context, cycleID string) (int, error) {
	var n int
	err := d.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE cycle_id=? AND state=?`, cycleID, model.AlertOpen).Scan(&n)
	return n, err
}
func (d *DB) AcknowledgeAlert(ctx context.Context, id, when string) error {
	_, err := d.Exec(ctx, `UPDATE alerts SET state=?,acknowledged_at=? WHERE id=?`, model.AlertAcknowledged, when, id)
	return err
}
