package model

func CloneVat(value Vat) Vat                      { return value }
func CloneCycle(value Cycle) Cycle                { return value }
func CloneSample(value Sample) Sample             { return value }
func CloneInspection(value Inspection) Inspection { return value }
func CloneAlert(value Alert) Alert                { return value }
func CloneAudit(value AuditEvent) AuditEvent      { return value }

func CloneRecipe(value Recipe) Recipe {
	out := value
	out.Stages = value.Stages
	return out
}

func CloneReport(value ReleaseReport) ReleaseReport {
	out := value
	out.Recipe = CloneRecipe(value.Recipe)
	out.Samples = cloneSamples(value.Samples)
	out.Inspections = make([]Inspection, len(value.Inspections))
	copy(out.Inspections, value.Inspections)
	out.Alerts = value.Alerts
	out.Audit = value.Audit
	return out
}

func cloneSamples(value []Sample) []Sample {
	out := make([]Sample, len(value))
	copy(out, value)
	return out
}
