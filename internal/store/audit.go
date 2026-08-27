package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) AppendAudit(ctx context.Context, e model.AuditEvent) error {
	_, err := d.Exec(ctx, `INSERT INTO audit_events(id,entity_type,entity_id,action,detail,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.EntityType, e.EntityID, e.Action, e.Detail, Text(e.CreatedAt))
	return err
}
func (d *DB) ListAudit(ctx context.Context, typ, id string) ([]model.AuditEvent, error) {
	rows, err := d.Query(ctx, `SELECT id,entity_type,entity_id,action,detail,created_at FROM audit_events WHERE entity_type=? AND entity_id=? ORDER BY created_at`, typ, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.AuditEvent{}
	for rows.Next() {
		var e model.AuditEvent
		var c string
		if err := rows.Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.Detail, &c); err != nil {
			return nil, err
		}
		e.CreatedAt = Time(c)
		out = append(out, e)
	}
	return out, rows.Err()
}
