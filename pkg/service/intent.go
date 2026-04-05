package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-verify-service/pkg/store"
	"github.com/google/uuid"
)

var IntentService = &intentService{}

type intentService struct{}

func (s *intentService) UpdateStatus(ctx context.Context, id uuid.UUID, status string, externalRef, message string) error {
	result, err := store.DB().ExecContext(ctx, `
		update execution_intents
		set status = $2, external_ref = $3, message = $4, last_error = '', updated_at = $5
		where id = $1 and deleted_at is null
	`, id, status, externalRef, message, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *intentService) UpdateStatusByResource(ctx context.Context, kind string, resourceID uuid.UUID, status string, externalRef, message string) error {
	result, err := store.DB().ExecContext(ctx, `
		update execution_intents
		set status = $3, external_ref = $4, message = $5, updated_at = $6
		where kind = $1 and resource_id = $2 and deleted_at is null
	`, kind, resourceID, status, externalRef, message, time.Now())
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}
