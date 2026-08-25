package model

func CloneVat(value Vat) Vat                      { return value }
func CloneCycle(value Cycle) Cycle                { return value }
func CloneSample(value Sample) Sample             { return value }
func CloneInspection(value Inspection) Inspection { return value }
func CloneAlert(value Alert) Alert                { return value }
func CloneAudit(value AuditEvent) AuditEvent      { return value }

func CloneRecipe(value Recipe) Recipe {
	out := value
	out.Stages = append([]RecipeStage(nil), value.Stages...)
	return out
}

func CloneReport(value ReleaseReport) ReleaseReport {
	out := value
	out.Recipe = CloneRecipe(value.Recipe)
	out.Samples = append([]Sample(nil), value.Samples...)
	out.Inspections = append([]Inspection(nil), value.Inspections...)
	out.Alerts = append([]Alert(nil), value.Alerts...)
	out.Audit = append([]AuditEvent(nil), value.Audit...)
	return out
}
