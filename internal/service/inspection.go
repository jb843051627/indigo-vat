package service

import (
	"context"

	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
)

func (s *Service) Inspect(ctx context.Context, in model.InspectionInput) (model.Inspection, error) {
	if err := validation.ValidateInspection(in); err != nil {
		return model.Inspection{}, err
	}
	now := s.clock.Now()
	i := model.Inspection{ID: validation.NewID("inspection"), CycleID: in.CycleID, Kind: in.Kind, Result: model.InspectionOpen, Score: in.Score, Inspector: in.Inspector, CreatedAt: now, Note: in.Note}
	if err := s.db.PutInspection(ctx, i); err != nil {
		return model.Inspection{}, err
	}
	return i, nil
}
func (s *Service) CompleteInspection(ctx context.Context, id string) (model.Inspection, error) {
	items, err := s.db.ListInspections(ctx, "")
	if err != nil {
		return model.Inspection{}, err
	}
	for _, item := range items {
		if item.ID == id {
			item.Result = model.InspectionPass
			if !validation.IsPassable(item.Score) {
				item.Result = model.InspectionFail
			}
			item.CompletedAt = s.clock.Now()
			if err := s.db.PutInspection(ctx, item); err != nil {
				return model.Inspection{}, err
			}
			return item, nil
		}
	}
	return model.Inspection{}, nil
}
func (s *Service) ListInspections(ctx context.Context, cycleID string) ([]model.Inspection, error) {
	return s.db.ListInspections(ctx, cycleID)
}
