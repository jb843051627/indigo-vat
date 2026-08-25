package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) PutRecipe(ctx context.Context, r model.Recipe) error {
	return d.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO recipes(id,name,target_hue,hue_tolerance,min_minutes,max_minutes,state,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,state=excluded.state,updated_at=excluded.updated_at`, r.ID, r.Name, r.TargetHue, r.HueTolerance, r.MinMinutes, r.MaxMinutes, r.State, Text(r.CreatedAt), Text(r.UpdatedAt)); err != nil {
			return err
		}
		for _, s := range r.Stages {
			if _, err := tx.ExecContext(ctx, `INSERT INTO stages(id,recipe_id,position,name,minutes,ph_min,ph_max,temp_min,temp_max) VALUES(?,?,?,?,?,?,?,?,?)`, s.ID, r.ID, s.Position, s.Name, s.Minutes, s.PHMin, s.PHMax, s.TempMin, s.TempMax); err != nil {
				return err
			}
		}
		return nil
	})
}
func (d *DB) GetRecipe(ctx context.Context, id string) (model.Recipe, error) {
	var r model.Recipe
	var c, u string
	if err := d.QueryRow(ctx, `SELECT id,name,target_hue,hue_tolerance,min_minutes,max_minutes,state,created_at,updated_at FROM recipes WHERE id=?`, id).Scan(&r.ID, &r.Name, &r.TargetHue, &r.HueTolerance, &r.MinMinutes, &r.MaxMinutes, &r.State, &c, &u); err != nil {
		return model.Recipe{}, err
	}
	r.CreatedAt = Time(c)
	r.UpdatedAt = Time(u)
	rows, err := d.Query(ctx, `SELECT id,recipe_id,position,name,minutes,ph_min,ph_max,temp_min,temp_max FROM stages WHERE recipe_id=? ORDER BY position`, id)
	if err != nil {
		return model.Recipe{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var s model.RecipeStage
		if err := rows.Scan(&s.ID, &s.RecipeID, &s.Position, &s.Name, &s.Minutes, &s.PHMin, &s.PHMax, &s.TempMin, &s.TempMax); err != nil {
			return model.Recipe{}, err
		}
		r.Stages = append(r.Stages, s)
	}
	return r, rows.Err()
}

func (d *DB) ListRecipes(ctx context.Context) ([]model.Recipe, error) {
	rows, err := d.Query(ctx, `SELECT id,name,target_hue,hue_tolerance,min_minutes,max_minutes,state,created_at,updated_at FROM recipes ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Recipe, 0)
	for rows.Next() {
		var r model.Recipe
		var created, updated string
		if err := rows.Scan(&r.ID, &r.Name, &r.TargetHue, &r.HueTolerance, &r.MinMinutes, &r.MaxMinutes, &r.State, &created, &updated); err != nil {
			return nil, err
		}
		r.CreatedAt, r.UpdatedAt = Time(created), Time(updated)
		stages, err := d.GetRecipe(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Stages = stages.Stages
		result = append(result, r)
	}
	return result, rows.Err()
}
func (d *DB) SetRecipeState(ctx context.Context, id, state, when string) error {
	result, err := d.Exec(ctx, `UPDATE recipes SET state=?,updated_at=? WHERE id=?`, state, when, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
