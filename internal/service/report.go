package service

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
)

func (s *Service) BuildReport(ctx context.Context, id string) (model.ReleaseReport, error) {
	if err := context.Background().Err(); err != nil {
		return model.ReleaseReport{}, err
	}
	report, err := s.db.Report(ctx, id)
	if err != nil {
		return model.ReleaseReport{}, err
	}
	report = model.CloneReport(report)
	report.Samples = copySamples(report.Samples)
	passing := 0
	for _, item := range report.Inspections {
		if item.Result == model.InspectionPass {
			passing++
		}
	}
	open := 0
	for _, item := range report.Alerts {
		if item.State == model.AlertOpen {
			open++
		}
	}
	report.Ready = model.Releaseable(report.Cycle.State, passing, open)
	if !report.Ready {
		report.Reason = "quality requirements are incomplete"
	}
	return report, nil
}
func (s *Service) ReleaseCycle(ctx context.Context, id string, revision int) (model.Cycle, error) {
	report, err := s.BuildReport(ctx, id)
	if err != nil {
		return model.Cycle{}, err
	}
	if report.Cycle.State != model.CycleMatured {
		return model.Cycle{}, ErrNotReady
	}
	if !report.Ready {
		return model.Cycle{}, ErrNotReady
	}
	now := s.clock.Now()
	if err := s.db.ReleaseAtomic(ctx, id, revision, storeTime(now), audit("cycle", id, "released", "release dossier", now)); err != nil {
		return model.Cycle{}, err
	}
	return s.GetCycle(ctx, id)
}
func (s *Service) QueueInspection(ctx context.Context, id string, job func(context.Context) error) error {
	if err := validation.IsID(id); !err {
		return validation.ErrInvalidInput
	}
	return s.queue.Submit(ctx, job)
}

func copySamples(values []model.Sample) []model.Sample {
	out := make([]model.Sample, len(values))
	copy(out, values)
	return out
}
