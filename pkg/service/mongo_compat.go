package service

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/google/uuid"
)

type releaseRecord struct {
	ID        uuid.UUID
	Status    model.ReleaseStatus
	Type      string
	Steps     []model.ReleaseStep
	DeletedAt *time.Time
}

type imageRecord struct {
	ID         uuid.UUID
	PipelineID string
	Status     model.ImageStatus
	Steps      []model.ImageStep
	DeletedAt  *time.Time
}

func scanReleaseRecord(scanner interface {
	Scan(dest ...any) error
}) (*releaseRecord, error) {
	var (
		record    releaseRecord
		stepsJSON []byte
		deletedAt sql.NullTime
	)
	if err := scanner.Scan(&record.ID, &record.Type, &stepsJSON, &record.Status, &deletedAt); err != nil {
		return nil, err
	}
	if len(stepsJSON) > 0 {
		if err := json.Unmarshal(stepsJSON, &record.Steps); err != nil {
			return nil, err
		}
	}
	if deletedAt.Valid {
		record.DeletedAt = &deletedAt.Time
	}
	return &record, nil
}

func scanImageRecord(scanner interface {
	Scan(dest ...any) error
}) (*imageRecord, error) {
	var (
		record    imageRecord
		stepsJSON []byte
		deletedAt sql.NullTime
	)
	if err := scanner.Scan(&record.ID, &record.PipelineID, &stepsJSON, &record.Status, &deletedAt); err != nil {
		return nil, err
	}
	if len(stepsJSON) > 0 {
		if err := json.Unmarshal(stepsJSON, &record.Steps); err != nil {
			return nil, err
		}
	}
	if deletedAt.Valid {
		record.DeletedAt = &deletedAt.Time
	}
	return &record, nil
}

func marshalJSON(value any, empty string) ([]byte, error) {
	if value == nil {
		return []byte(empty), nil
	}
	return json.Marshal(value)
}

func nullableTimePtr(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
}
