package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutCycle(ctx context.Context, c model.Cycle) error {
	_, err := d.Exec(ctx, `INSERT INTO cycles(id,vat_id,recipe_id,state,revision,started_at,matured_at,released_at,updated_at,note) VALUES(?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET state=excluded.state,revision=excluded.revision,matured_at=excluded.matured_at,released_at=excluded.released_at,updated_at=excluded.updated_at,note=excluded.note`, c.ID, c.VatID, c.RecipeID, c.State, c.Revision, Text(c.StartedAt), Text(c.MaturedAt), Text(c.ReleasedAt), Text(c.UpdatedAt), c.Note)
	return err
}
func (d *DB) GetCycle(ctx context.Context, id string) (model.Cycle, error) {
	var c model.Cycle
	var a, b, e, u string
	err := d.QueryRow(ctx, `SELECT id,vat_id,recipe_id,state,revision,started_at,matured_at,released_at,updated_at,note FROM cycles WHERE id=?`, id).Scan(&c.ID, &c.VatID, &c.RecipeID, &c.State, &c.Revision, &a, &b, &e, &u, &c.Note)
	c.StartedAt = Time(a)
	c.MaturedAt = Time(b)
	c.ReleasedAt = Time(e)
	c.UpdatedAt = Time(u)
	return c, err
}
func (d *DB) ListCycles(ctx context.Context, state string) ([]model.Cycle, error) {
	q := `SELECT id,vat_id,recipe_id,state,revision,started_at,matured_at,released_at,updated_at,note FROM cycles`
	args := []any{}
	if state != "" {
		q += ` WHERE state=?`
		args = append(args, state)
	}
	rows, err := d.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Cycle{}
	for rows.Next() {
		var c model.Cycle
		var a, b, e, u string
		if err := rows.Scan(&c.ID, &c.VatID, &c.RecipeID, &c.State, &c.Revision, &a, &b, &e, &u, &c.Note); err != nil {
			return nil, err
		}
		c.StartedAt = Time(a)
		c.MaturedAt = Time(b)
		c.ReleasedAt = Time(e)
		c.UpdatedAt = Time(u)
		out = append(out, c)
	}
	return out, rows.Err()
}
func (d *DB) UpdateCycle(ctx context.Context, id, state, note, when string, revision int, matured, released string) error {
	result, err := d.Exec(ctx, `UPDATE cycles SET state=?,revision=revision+1,note=?,updated_at=?,matured_at=?,released_at=? WHERE id=? AND revision=?`, state, note, when, matured, released, id, revision)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
