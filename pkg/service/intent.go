package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var IntentService = &intentService{}

type intentService struct{}

func (s *intentService) UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.IntentStatus, externalRef, message string) error {
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"external_ref": externalRef,
			"message":      message,
			"last_error":   "",
			"updated_at":   time.Now(),
		},
	}

	return mongo.Repo.UpdateByID(ctx, &model.Intent{}, id, update)
}

func (s *intentService) UpdateStatusByResource(ctx context.Context, kind model.IntentKind, resourceID primitive.ObjectID, status model.IntentStatus, externalRef, message string) error {
	filter := bson.M{
		"kind":        kind,
		"resource_id": resourceID,
	}
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"external_ref": externalRef,
			"message":      message,
			"updated_at":   time.Now(),
		},
	}

	return mongo.Repo.UpdateOne(ctx, &model.Intent{}, filter, update)
}
