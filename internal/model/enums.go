package model

const (
	VatActive         = "active"
	VatPaused         = "paused"
	VatRetired        = "retired"
	RecipeDraft       = "draft"
	RecipeReady       = "ready"
	CyclePlanned      = "planned"
	CycleFermenting   = "fermenting"
	CycleMatured      = "matured"
	CycleHeld         = "held"
	CycleReleased     = "released"
	CycleDiscarded    = "discarded"
	SamplePending     = "pending"
	SampleAccepted    = "accepted"
	SampleRejected    = "rejected"
	InspectionOpen    = "open"
	InspectionPass    = "pass"
	InspectionFail    = "fail"
	AlertOpen         = "open"
	AlertAcknowledged = "acknowledged"
	AlertInfo         = "info"
	AlertWarn         = "warn"
	AlertCritical     = "critical"
)

func AllowedCycleState(value string) bool {
	return value == CyclePlanned || value == CycleFermenting || value == CycleMatured || value == CycleHeld || value == CycleReleased || value == CycleDiscarded
}

func CanTransition(from, to string) bool {
	switch from {
	case CyclePlanned:
		return to == CycleFermenting || to == CycleDiscarded
	case CycleFermenting:
		return to == CycleMatured || to == CycleHeld || to == CycleDiscarded
	case CycleMatured:
		return to == CycleHeld || to == CycleReleased || to == CycleDiscarded
	case CycleHeld:
		return to == CycleFermenting || to == CycleMatured || to == CycleDiscarded
	case CycleReleased:
		return to == CycleFermenting
	case CycleDiscarded:
		return false
	default:
		return false
	}
}

func IsTerminal(state string) bool { return state == CycleReleased || state == CycleDiscarded }

func IsQualityState(state string) bool {
	return state == CycleMatured || state == CycleHeld || state == CycleReleased
}

func Releaseable(state string, passing, openAlerts int) bool {
	return state == CycleMatured && passing > 0 && openAlerts == 0
}
