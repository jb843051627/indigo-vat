package model

type VatInput struct {
	Name       string  `json:"name"`
	Site       string  `json:"site"`
	Capacity   int     `json:"capacity"`
	TargetPH   float64 `json:"target_ph"`
	TargetTemp float64 `json:"target_temp"`
	Timezone   string  `json:"timezone"`
}
type RecipeInput struct {
	Name         string  `json:"name"`
	TargetHue    float64 `json:"target_hue"`
	HueTolerance float64 `json:"hue_tolerance"`
	MinMinutes   int     `json:"min_minutes"`
	MaxMinutes   int     `json:"max_minutes"`
}
type StageInput struct {
	Name    string  `json:"name"`
	Minutes int     `json:"minutes"`
	PHMin   float64 `json:"ph_min"`
	PHMax   float64 `json:"ph_max"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}
type CycleInput struct {
	VatID    string `json:"vat_id"`
	RecipeID string `json:"recipe_id"`
	Note     string `json:"note"`
}
type SampleInput struct {
	CycleID     string  `json:"cycle_id"`
	TakenAt     string  `json:"taken_at"`
	Hue         float64 `json:"hue"`
	PH          float64 `json:"ph"`
	Temperature float64 `json:"temperature"`
	Observer    string  `json:"observer"`
	Note        string  `json:"note"`
}
type InspectionInput struct {
	CycleID   string  `json:"cycle_id"`
	Kind      string  `json:"kind"`
	Inspector string  `json:"inspector"`
	Score     float64 `json:"score"`
	Note      string  `json:"note"`
}
type AlertInput struct {
	CycleID string `json:"cycle_id"`
	Level   string `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
type TransitionInput struct {
	State            string `json:"state"`
	ExpectedRevision int    `json:"expected_revision"`
	Note             string `json:"note"`
}
