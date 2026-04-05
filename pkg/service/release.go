package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/bsonger/devflow-verify-service/pkg/store"
	"github.com/google/uuid"
)

var ReleaseService = &releaseService{}

type releaseService struct{}

func (s *releaseService) Get(ctx context.Context, id uuid.UUID) (*releaseRecord, error) {
	record, err := scanReleaseRecord(store.DB().QueryRowContext(ctx, `
		select id, type, steps, status, deleted_at
		from releases
		where id = $1
	`, id))
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, sql.ErrNoRows
	}
	return record, nil
}

func (s *releaseService) updateStatus(ctx context.Context, releaseID uuid.UUID, status model.ReleaseStatus) error {
	record, err := s.Get(ctx, releaseID)
	if err != nil {
		return err
	}
	switch record.Status {
	case model.ReleaseSucceeded, model.ReleaseFailed, model.ReleaseRolledBack, model.ReleaseSyncFailed:
		return nil
	}
	if record.Status == status {
		return nil
	}
	result, err := store.DB().ExecContext(ctx, `
		update releases
		set status = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, releaseID, status, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *releaseService) UpdateStatus(ctx context.Context, releaseID uuid.UUID, status model.ReleaseStatus) error {
	return s.updateStatus(ctx, releaseID, status)
}

func (s *releaseService) UpdateStep(ctx context.Context, releaseID uuid.UUID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	if stepName == "" {
		return nil
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	release, err := s.Get(ctx, releaseID)
	if err != nil {
		return err
	}

	nextSteps := cloneReleaseSteps(release.Steps)
	currentStep := findReleaseStep(release.Steps, stepName)
	if currentStep == nil {
		nextSteps = append(nextSteps, model.ReleaseStep{Name: stepName, Progress: progress, Status: status, Message: message, StartTime: start, EndTime: end})
		if err := s.updateSteps(ctx, releaseID, nextSteps); err != nil {
			return err
		}
		return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
	}

	if currentStep.Status == model.StepFailed || currentStep.Status == model.StepSucceeded {
		return nil
	}

	applyReleaseStepUpdate(nextSteps, stepName, status, progress, message, start, end)
	if err := s.updateSteps(ctx, releaseID, nextSteps); err != nil {
		return err
	}
	return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
}

func findReleaseStep(steps []model.ReleaseStep, stepName string) *model.ReleaseStep {
	for _, step := range steps {
		if step.Name == stepName {
			current := step
			return &current
		}
	}
	return nil
}

func cloneReleaseSteps(steps []model.ReleaseStep) []model.ReleaseStep {
	if len(steps) == 0 {
		return nil
	}
	cloned := make([]model.ReleaseStep, len(steps))
	copy(cloned, steps)
	return cloned
}

func applyReleaseStepUpdate(steps []model.ReleaseStep, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) {
	for i := range steps {
		if steps[i].Name != stepName {
			continue
		}
		steps[i].Status = status
		steps[i].Progress = progress
		steps[i].Message = message
		if start != nil {
			steps[i].StartTime = start
		}
		if end != nil {
			steps[i].EndTime = end
		}
		return
	}
}

func (s *releaseService) updateStatusFromSteps(ctx context.Context, releaseID uuid.UUID, releaseAction string, currentStatus model.ReleaseStatus, steps []model.ReleaseStep) error {
	nextStatus := model.DeriveReleaseStatusFromSteps(releaseAction, currentStatus, steps)
	if nextStatus == currentStatus {
		return nil
	}
	return s.updateStatus(ctx, releaseID, nextStatus)
}

func (s *releaseService) updateSteps(ctx context.Context, releaseID uuid.UUID, steps []model.ReleaseStep) error {
	stepsJSON, err := marshalJSON(steps, "[]")
	if err != nil {
		return err
	}
	result, err := store.DB().ExecContext(ctx, `
		update releases
		set steps = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, releaseID, stepsJSON, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}
