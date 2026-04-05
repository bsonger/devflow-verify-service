package service

import (
	"context"
	"errors"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var ManifestService = &manifestService{}

type manifestService struct{}

func (s *manifestService) Get(ctx context.Context, id uuid.UUID) (*manifestRecord, error) {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return nil, err
	}
	doc := &manifestDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, oid); err != nil {
		return nil, err
	}
	record := manifestRecordFromDoc(doc)
	return &record, nil
}

func (s *manifestService) AssignPipelineID(ctx context.Context, manifestID uuid.UUID, pipelineID string) error {
	if manifestID == uuid.Nil {
		return errors.New("manifest id cannot be zero")
	}
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	oid, err := bridgeUUIDToObjectID(manifestID)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateByID(ctx, &manifestDoc{}, oid, bson.M{
		"$set": bson.M{"pipeline_id": pipelineID, "updated_at": time.Now()},
	})
}

func (s *manifestService) UpdateManifestStatusByID(ctx context.Context, manifestID uuid.UUID, status model.ManifestStatus) error {
	if manifestID == uuid.Nil {
		return errors.New("manifest id cannot be zero")
	}
	oid, err := bridgeUUIDToObjectID(manifestID)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateByID(ctx, &manifestDoc{}, oid, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now()},
	})
}

func (s *manifestService) UpdateStepStatus(ctx context.Context, pipelineID, taskName string, status model.StepStatus, message string, start, end *time.Time) error {
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	if taskName == "" {
		return errors.New("task name cannot be empty")
	}

	update := bson.M{"steps.$.status": status, "steps.$.message": message, "updated_at": time.Now()}
	if start != nil {
		update["steps.$.start_time"] = *start
	}
	if end != nil {
		update["steps.$.end_time"] = *end
	}

	return mongo.Repo.UpdateOne(ctx, &manifestDoc{}, bson.M{
		"pipeline_id": pipelineID,
		"steps": bson.M{"$elemMatch": bson.M{
			"task_name": taskName,
			"status":    bson.M{"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded, status}},
		}},
	}, bson.M{"$set": update})
}

func (s *manifestService) BindTaskRun(ctx context.Context, pipelineID, taskName, taskRun string) error {
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	if taskName == "" {
		return errors.New("task name cannot be empty")
	}
	if taskRun == "" {
		return errors.New("task run cannot be empty")
	}

	return mongo.Repo.UpdateOne(ctx, &manifestDoc{}, bson.M{
		"pipeline_id": pipelineID,
		"steps": bson.M{"$elemMatch": bson.M{
			"task_name": taskName,
			"task_run":  bson.M{"$exists": false},
			"status":    bson.M{"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded}},
		}},
	}, bson.M{"$set": bson.M{"steps.$.task_run": taskRun, "updated_at": time.Now()}})
}
