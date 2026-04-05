package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var IntentService = &intentService{}

type intentService struct{}

func (s *intentService) UpdateStatus(ctx context.Context, id uuid.UUID, status string, externalRef, message string) error {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateByID(ctx, &intentDoc{}, oid, bson.M{
		"$set": bson.M{
			"status":       status,
			"external_ref": externalRef,
			"message":      message,
			"last_error":   "",
			"updated_at":   time.Now(),
		},
	})
}

func (s *intentService) UpdateStatusByResource(ctx context.Context, kind string, resourceID uuid.UUID, status string, externalRef, message string) error {
	resourceOID, err := bridgeUUIDToObjectID(resourceID)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateOne(ctx, &intentDoc{}, bson.M{
		"kind":        kind,
		"resource_id": resourceOID,
	}, bson.M{
		"$set": bson.M{
			"status":       status,
			"external_ref": externalRef,
			"message":      message,
			"updated_at":   time.Now(),
		},
	})
}
