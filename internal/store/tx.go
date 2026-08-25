package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) CreateCycleAtomic(ctx context.Context, c model.Cycle, e model.AuditEvent) error {
	return d.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO cycles(id,vat_id,recipe_id,state,revision,started_at,matured_at,released_at,updated_at,note) VALUES(?,?,?,?,?,?,?,?,?,?)`, c.ID, c.VatID, c.RecipeID, c.State, c.Revision, Text(c.StartedAt), Text(c.MaturedAt), Text(c.ReleasedAt), Text(c.UpdatedAt), c.Note); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, `INSERT INTO audit_events(id,entity_type,entity_id,action,detail,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.EntityType, e.EntityID, e.Action, e.Detail, Text(e.CreatedAt))
		return err
	})
}
func (d *DB) CreateSampleAtomic(ctx context.Context, s model.Sample, e model.AuditEvent) error {
	if _, err := d.Exec(ctx, `INSERT INTO samples(id,cycle_id,taken_at,hue,ph,temperature,status,observer,note) VALUES(?,?,?,?,?,?,?,?,?)`, s.ID, s.CycleID, Text(s.TakenAt), s.Hue, s.PH, s.Temperature, s.Status, s.Observer, s.Note); err != nil {
		return err
	}
	_, err := d.Exec(ctx, `INSERT INTO audit_events(id,entity_type,entity_id,action,detail,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.EntityType, e.EntityID, e.Action, e.Detail, Text(e.CreatedAt))
	return err
}
func (d *DB) ReleaseAtomic(ctx context.Context, id string, revision int, when string, e model.AuditEvent) error {
	return d.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `UPDATE cycles SET state=?,revision=revision+1,released_at=?,updated_at=? WHERE id=? AND revision=?`, model.CycleReleased, when, when, id, revision)
		if err != nil {
			return err
		}
		n, _ := result.RowsAffected()
		if n == 0 {
			return sql.ErrNoRows
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO audit_events(id,entity_type,entity_id,action,detail,created_at) VALUES(?,?,?,?,?,?)`, e.ID, e.EntityType, e.EntityID, e.Action, e.Detail, Text(e.CreatedAt))
		return err
	})
}
