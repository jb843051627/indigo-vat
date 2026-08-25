package model

import "time"

type Vat struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Site       string    `json:"site"`
	State      string    `json:"state"`
	Capacity   int       `json:"capacity"`
	TargetPH   float64   `json:"target_ph"`
	TargetTemp float64   `json:"target_temp"`
	Timezone   string    `json:"timezone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Recipe struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	TargetHue    float64       `json:"target_hue"`
	HueTolerance float64       `json:"hue_tolerance"`
	MinMinutes   int           `json:"min_minutes"`
	MaxMinutes   int           `json:"max_minutes"`
	State        string        `json:"state"`
	Stages       []RecipeStage `json:"stages"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type RecipeStage struct {
	ID       string  `json:"id"`
	RecipeID string  `json:"recipe_id"`
	Position int     `json:"position"`
	Name     string  `json:"name"`
	Minutes  int     `json:"minutes"`
	PHMin    float64 `json:"ph_min"`
	PHMax    float64 `json:"ph_max"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

type Cycle struct {
	ID         string    `json:"id"`
	VatID      string    `json:"vat_id"`
	RecipeID   string    `json:"recipe_id"`
	State      string    `json:"state"`
	Revision   int       `json:"revision"`
	StartedAt  time.Time `json:"started_at"`
	MaturedAt  time.Time `json:"matured_at"`
	ReleasedAt time.Time `json:"released_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Note       string    `json:"note"`
}

type Sample struct {
	ID          string    `json:"id"`
	CycleID     string    `json:"cycle_id"`
	TakenAt     time.Time `json:"taken_at"`
	Hue         float64   `json:"hue"`
	PH          float64   `json:"ph"`
	Temperature float64   `json:"temperature"`
	Status      string    `json:"status"`
	Observer    string    `json:"observer"`
	Note        string    `json:"note"`
}

type Inspection struct {
	ID          string    `json:"id"`
	CycleID     string    `json:"cycle_id"`
	Kind        string    `json:"kind"`
	Result      string    `json:"result"`
	Score       float64   `json:"score"`
	Inspector   string    `json:"inspector"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Note        string    `json:"note"`
}

type Alert struct {
	ID             string    `json:"id"`
	CycleID        string    `json:"cycle_id"`
	Level          string    `json:"level"`
	Code           string    `json:"code"`
	Message        string    `json:"message"`
	State          string    `json:"state"`
	CreatedAt      time.Time `json:"created_at"`
	AcknowledgedAt time.Time `json:"acknowledged_at"`
}

type AuditEvent struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReleaseReport struct {
	Cycle       Cycle        `json:"cycle"`
	Vat         Vat          `json:"vat"`
	Recipe      Recipe       `json:"recipe"`
	Samples     []Sample     `json:"samples"`
	Inspections []Inspection `json:"inspections"`
	Alerts      []Alert      `json:"alerts"`
	Audit       []AuditEvent `json:"audit"`
	Ready       bool         `json:"ready"`
	Reason      string       `json:"reason"`
}
