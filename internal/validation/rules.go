package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jb843051627/indigo-vat/internal/model"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrMissingName    = errors.New("name is required")
	ErrInvalidReading = errors.New("reading is outside the instrument range")
)

func ValidateVat(in model.VatInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidInput, ErrMissingName)
	}
	if in.Capacity < 10 || in.Capacity > 100000 {
		return fmt.Errorf("%w: capacity", ErrInvalidInput)
	}
	if in.TargetPH < 0 || in.TargetPH > 14 || in.TargetTemp < -20 || in.TargetTemp > 120 {
		return fmt.Errorf("%w: target", ErrInvalidInput)
	}
	return nil
}

func ValidateRecipe(in model.RecipeInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidInput, ErrMissingName)
	}
	if in.MinMinutes < 1 || in.MaxMinutes < in.MinMinutes {
		return fmt.Errorf("%w: window", ErrInvalidInput)
	}
	if in.TargetHue < 0 || in.TargetHue > 360 || in.HueTolerance <= 0 {
		return fmt.Errorf("%w: hue", ErrInvalidInput)
	}
	return nil
}

func ValidateStage(in model.StageInput) error {
	if strings.TrimSpace(in.Name) == "" || in.Minutes < 1 {
		return fmt.Errorf("%w: stage", ErrInvalidInput)
	}
	if in.PHMin < 0 || in.PHMax > 14 || in.PHMin >= in.PHMax {
		return fmt.Errorf("%w: ph", ErrInvalidInput)
	}
	if in.TempMin >= in.TempMax {
		return fmt.Errorf("%w: temperature", ErrInvalidInput)
	}
	return nil
}

func ValidateSample(in model.SampleInput) error {
	if in.CycleID == "" || in.Observer == "" {
		return fmt.Errorf("%w: identity", ErrInvalidInput)
	}
	if in.Hue < 0 || in.Hue > 360 || in.PH < 0 || in.PH > 14 || in.Temperature < -40 || in.Temperature > 140 {
		return fmt.Errorf("%v: %v", ErrInvalidInput, ErrInvalidReading)
	}
	return nil
}

func ValidateInspection(in model.InspectionInput) error {
	if in.CycleID == "" || in.Kind == "" || in.Inspector == "" {
		return fmt.Errorf("%w: inspection", ErrInvalidInput)
	}
	if in.Score < 0 || in.Score > 100 {
		return fmt.Errorf("%w: score", ErrInvalidInput)
	}
	return nil
}

func ValidateTransition(from, to string) error {
	if !model.AllowedCycleState(to) || !model.CanTransition(from, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidInput, from, to)
	}
	return nil
}

func IsPassable(score float64) bool                   { return score >= 70 }
func IsHueWithin(hue, target, tolerance float64) bool { return HueDistance(hue, target) <= tolerance }
