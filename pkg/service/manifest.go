package service

import (
	"context"
	"errors"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ManifestService = &manifestService{}

type manifestService struct{}

func (s *manifestService) Get(ctx context.Context, id primitive.ObjectID) (*model.Manifest, error) {
	manifest := &model.Manifest{}
	if err := mongo.Repo.FindByID(ctx, manifest, id); err != nil {
		return nil, err
	}
	return manifest, nil
}

func (s *manifestService) AssignPipelineID(ctx context.Context, manifestID primitive.ObjectID, pipelineID string) error {
	if manifestID.IsZero() {
		return errors.New("manifest id cannot be zero")
	}
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}

	return mongo.Repo.UpdateByID(ctx, &model.Manifest{}, manifestID, bson.M{
		"$set": bson.M{
			"pipeline_id": pipelineID,
			"updated_at":  time.Now(),
		},
	})
}

func (s *manifestService) UpdateManifestStatusByID(ctx context.Context, manifestID primitive.ObjectID, status model.ManifestStatus) error {
	if manifestID.IsZero() {
		return errors.New("manifest id cannot be zero")
	}

	return mongo.Repo.UpdateByID(ctx, &model.Manifest{}, manifestID, bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	})
}

func (s *manifestService) UpdateStepStatus(ctx context.Context, pipelineID, taskName string, status model.StepStatus, message string, start, end *time.Time) error {
	if pipelineID == "" {
		return errors.New("pipeline id cannot be empty")
	}
	if taskName == "" {
		return errors.New("task name cannot be empty")
	}

	update := bson.M{
		"steps.$.status":  status,
		"steps.$.message": message,
		"updated_at":      time.Now(),
	}

	if start != nil {
		update["steps.$.start_time"] = *start
	}
	if end != nil {
		update["steps.$.end_time"] = *end
	}

	filter := bson.M{
		"pipeline_id": pipelineID,
		"steps": bson.M{
			"$elemMatch": bson.M{
				"task_name": taskName,
				"status": bson.M{
					"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded, status},
				},
			},
		},
	}

	return mongo.Repo.UpdateOne(ctx, &model.Manifest{}, filter, bson.M{"$set": update})
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

	return mongo.Repo.UpdateOne(
		ctx,
		&model.Manifest{},
		bson.M{
			"pipeline_id": pipelineID,
			"steps": bson.M{
				"$elemMatch": bson.M{
					"task_name": taskName,
					"task_run":  bson.M{"$exists": false},
					"status": bson.M{
						"$nin": []model.StepStatus{
							model.StepFailed,
							model.StepSucceeded,
						},
					},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"steps.$.task_run": taskRun,
				"updated_at":       time.Now(),
			},
		},
	)
}
