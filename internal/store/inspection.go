package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutInspection(ctx context.Context, i model.Inspection) error {
	_, err := d.Exec(ctx, `INSERT INTO inspections(id,cycle_id,kind,result,score,inspector,created_at,completed_at,note) VALUES(?,?,?,?,?,?,?,?,?)`, i.ID, i.CycleID, i.Kind, i.Result, i.Score, i.Inspector, Text(i.CreatedAt), Text(i.CompletedAt), i.Note)
	return err
}
func (d *DB) ListInspections(ctx context.Context, cycleID string) ([]model.Inspection, error) {
	rows, err := d.Query(ctx, `SELECT id,cycle_id,kind,result,score,inspector,created_at,completed_at,note FROM inspections WHERE cycle_id=? ORDER BY created_at`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Inspection{}
	for rows.Next() {
		var i model.Inspection
		var a, b string
		if err := rows.Scan(&i.ID, &i.CycleID, &i.Kind, &i.Result, &i.Score, &i.Inspector, &a, &b, &i.Note); err != nil {
			return nil, err
		}
		i.CreatedAt = Time(a)
		i.CompletedAt = Time(b)
		out = append(out, i)
	}
	return out, rows.Err()
}
func (d *DB) PassingInspections(ctx context.Context, cycleID string) (int, error) {
	var n int
	err := d.QueryRow(ctx, `SELECT COUNT(*) FROM inspections WHERE cycle_id=? AND result=?`, cycleID, model.InspectionPass).Scan(&n)
	return n, err
}
