package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) Report(ctx context.Context, id string) (model.ReleaseReport, error) {
	c, err := d.GetCycle(ctx, id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	v, err := d.GetVat(ctx, c.VatID)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	r, err := d.GetRecipe(ctx, c.RecipeID)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	samples, err := d.ListSamples(ctx, id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	inspections, err := d.ListInspections(ctx, id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	alerts, err := d.ListAlerts(ctx, id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	audit, err := d.ListAudit(ctx, "cycle", id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	return model.ReleaseReport{Cycle: c, Vat: v, Recipe: r, Samples: samples, Inspections: inspections, Alerts: alerts, Audit: audit}, nil
}
