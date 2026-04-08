package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/bsonger/devflow-verify-service/pkg/store"
	"github.com/google/uuid"
)

var ImageService = &imageService{}

type imageService struct{}

func (s *imageService) Get(ctx context.Context, id uuid.UUID) (*imageRecord, error) {
	return scanImageRecord(store.DB().QueryRowContext(ctx, `
		select id, pipeline_id, steps, status, deleted_at
		from images
		where id = $1
	`, id))
}

func (s *imageService) AssignPipelineID(ctx context.Context, imageID uuid.UUID, pipelineID string) error {
	if imageID == uuid.Nil {
		return errors.New("image id cannot be zero")
	}
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	result, err := store.DB().ExecContext(ctx, `
		update images
		set pipeline_id = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, imageID, pipelineID, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *imageService) UpdateImageStatusByID(ctx context.Context, imageID uuid.UUID, status model.ImageStatus) error {
	if imageID == uuid.Nil {
		return errors.New("image id cannot be zero")
	}
	result, err := store.DB().ExecContext(ctx, `
		update images
		set status = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, imageID, status, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *imageService) UpdateStepStatus(ctx context.Context, pipelineID, taskName string, status model.StepStatus, message string, start, end *time.Time) error {
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	if taskName == "" {
		return errors.New("task name cannot be empty")
	}

	row := store.DB().QueryRowContext(ctx, `
		select id, pipeline_id, steps, status, deleted_at
		from images
		where pipeline_id = $1 and deleted_at is null
	`, pipelineID)
	record, err := scanImageRecord(row)
	if err != nil {
		return err
	}

	changed := false
	for i := range record.Steps {
		if record.Steps[i].TaskName != taskName {
			continue
		}
		if record.Steps[i].Status == model.StepFailed || record.Steps[i].Status == model.StepSucceeded || record.Steps[i].Status == status {
			return nil
		}
		record.Steps[i].Status = status
		record.Steps[i].Message = message
		if start != nil {
			record.Steps[i].StartTime = start
		}
		if end != nil {
			record.Steps[i].EndTime = end
		}
		changed = true
		break
	}
	if !changed {
		return nil
	}

	stepsJSON, err := marshalJSON(record.Steps, "[]")
	if err != nil {
		return err
	}
	result, err := store.DB().ExecContext(ctx, `
		update images
		set steps = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, record.ID, stepsJSON, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *imageService) BindTaskRun(ctx context.Context, pipelineID, taskName, taskRun string) error {
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	if taskName == "" {
		return errors.New("task name cannot be empty")
	}
	if taskRun == "" {
		return errors.New("task run cannot be empty")
	}

	row := store.DB().QueryRowContext(ctx, `
		select id, pipeline_id, steps, status, deleted_at
		from images
		where pipeline_id = $1 and deleted_at is null
	`, pipelineID)
	record, err := scanImageRecord(row)
	if err != nil {
		return err
	}

	changed := false
	for i := range record.Steps {
		if record.Steps[i].TaskName != taskName {
			continue
		}
		if record.Steps[i].TaskRun != "" || record.Steps[i].Status == model.StepFailed || record.Steps[i].Status == model.StepSucceeded {
			return nil
		}
		record.Steps[i].TaskRun = taskRun
		changed = true
		break
	}
	if !changed {
		return nil
	}

	stepsJSON, err := marshalJSON(record.Steps, "[]")
	if err != nil {
		return err
	}
	result, err := store.DB().ExecContext(ctx, `
		update images
		set steps = $2, updated_at = $3
		where id = $1 and deleted_at is null
	`, record.ID, stepsJSON, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func ensureRowsAffected(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
