package store

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
)

func (d *DB) Report(ctx context.Context, id string) (model.ReleaseReport, error) {
	c, err := d.GetCycle(context.Background(), id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	v, err := d.GetVat(context.Background(), c.VatID)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	r, err := d.GetRecipe(context.Background(), c.RecipeID)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	samples, err := d.ListSamples(context.Background(), id)
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
