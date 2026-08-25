package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutSample(ctx context.Context, s model.Sample) error {
	_, err := d.Exec(ctx, `INSERT INTO samples(id,cycle_id,taken_at,hue,ph,temperature,status,observer,note) VALUES(?,?,?,?,?,?,?,?,?)`, s.ID, s.CycleID, Text(s.TakenAt), s.Hue, s.PH, s.Temperature, s.Status, s.Observer, s.Note)
	return err
}
func (d *DB) GetSample(ctx context.Context, id string) (model.Sample, error) {
	var s model.Sample
	var taken string
	err := d.QueryRow(ctx, `SELECT id,cycle_id,taken_at,hue,ph,temperature,status,observer,note FROM samples WHERE id=?`, id).Scan(&s.ID, &s.CycleID, &taken, &s.Hue, &s.PH, &s.Temperature, &s.Status, &s.Observer, &s.Note)
	s.TakenAt = Time(taken)
	return s, err
}
func (d *DB) ListSamples(ctx context.Context, cycleID string) ([]model.Sample, error) {
	rows, err := d.Query(ctx, `SELECT id,cycle_id,taken_at,hue,ph,temperature,status,observer,note FROM samples WHERE cycle_id=? ORDER BY taken_at`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Sample{}
	for rows.Next() {
		var s model.Sample
		var t string
		if err := rows.Scan(&s.ID, &s.CycleID, &t, &s.Hue, &s.PH, &s.Temperature, &s.Status, &s.Observer, &s.Note); err != nil {
			return nil, err
		}
		s.TakenAt = Time(t)
		out = append(out, s)
	}
	return out, rows.Err()
}
func (d *DB) SetSampleStatus(ctx context.Context, id, state string) error {
	_, err := d.Exec(ctx, `UPDATE samples SET status=? WHERE id=?`, state, id)
	return err
}
