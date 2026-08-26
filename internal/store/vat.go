package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutVat(ctx context.Context, v model.Vat) error {
	_, err := d.Exec(ctx, `INSERT INTO vats(id,name,site,state,capacity,target_ph,target_temp,timezone,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,site=excluded.site,state=excluded.state,capacity=excluded.capacity,target_ph=excluded.target_ph,target_temp=excluded.target_temp,timezone=excluded.timezone,updated_at=excluded.updated_at`, v.ID, v.Name, v.Site, v.State, v.Capacity, v.TargetPH, v.TargetTemp, v.Timezone, Text(v.CreatedAt), Text(v.UpdatedAt))
	return err
}
func (d *DB) GetVat(ctx context.Context, id string) (model.Vat, error) {
	var v model.Vat
	var created, updated string
	err := d.QueryRow(ctx, `SELECT id,name,site,state,capacity,target_ph,target_temp,timezone,created_at,updated_at FROM vats WHERE id=?`, id).Scan(&v.ID, &v.Name, &v.Site, &v.State, &v.Capacity, &v.TargetPH, &v.TargetTemp, &v.Timezone, &created, &updated)
	v.CreatedAt = Time(created)
	v.UpdatedAt = Time(updated)
	return v, err
}
func (d *DB) ListVats(ctx context.Context) ([]model.Vat, error) {
	rows, err := d.Query(ctx, `SELECT id,name,site,state,capacity,target_ph,target_temp,timezone,created_at,updated_at FROM vats ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Vat{}
	for rows.Next() {
		var v model.Vat
		var c, u string
		if err := rows.Scan(&v.ID, &v.Name, &v.Site, &v.State, &v.Capacity, &v.TargetPH, &v.TargetTemp, &v.Timezone, &c, &u); err != nil {
			return nil, err
		}
		v.CreatedAt = Time(c)
		v.UpdatedAt = Time(u)
		out = append(out, v)
	}
	return out, rows.Err()
}
func (d *DB) SetVatState(ctx context.Context, id, state, when string) error {
	result, err := d.Exec(ctx, `UPDATE vats SET state=?,updated_at=? WHERE id=?`, state, when, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
