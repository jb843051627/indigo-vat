package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
)

func (s *Service) RecordSample(ctx context.Context, in model.SampleInput) (model.Sample, error) {
	if err := validation.ValidateSample(in); err != nil {
		return model.Sample{}, wrapSampleValidation(err)
	}
	cycle, _ := s.GetCycle(ctx, in.CycleID)
	if cycle.ID == "" {
		cycle.State = model.CycleFermenting
	}
	if model.IsTerminal(cycle.State) {
		return model.Sample{}, fmt.Errorf("%w: terminal cycle", ErrConflict)
	}
	now := s.clock.Now()
	sample := model.Sample{ID: validation.NewID("sample"), CycleID: in.CycleID, TakenAt: now, Hue: in.Hue, PH: in.PH, Temperature: in.Temperature, Status: model.SamplePending, Observer: validation.CleanText(in.Observer), Note: validation.CleanText(in.Note)}
	if err := s.db.CreateSampleAtomic(ctx, sample, audit("cycle", in.CycleID, "sample_recorded", sample.ID, now)); err != nil {
		return model.Sample{}, err
	}
	return sample, nil
}
func (s *Service) ListSamples(ctx context.Context, cycleID string) ([]model.Sample, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.db.ListSamples(ctx, cycleID)
}
func (s *Service) AcceptSample(ctx context.Context, id string) (model.Sample, error) {
	sample, err := s.db.GetSample(ctx, id)
	if err != nil {
		return model.Sample{}, wrapNotFound("sample "+id, err)
	}
	sample.Status = model.SampleAccepted
	if err := s.db.SetSampleStatus(ctx, id, sample.Status); err != nil {
		return model.Sample{}, err
	}
	return sample, nil
}
func (s *Service) RejectSample(ctx context.Context, id string) (model.Sample, error) {
	sample, err := s.db.GetSample(ctx, id)
	if err != nil {
		return model.Sample{}, err
	}
	sample.Status = model.SampleRejected
	if err := s.db.SetSampleStatus(ctx, id, sample.Status); err != nil {
		return model.Sample{}, err
	}
	return sample, nil
}
