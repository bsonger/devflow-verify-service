package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

var ReleaseService = &releaseService{}

type releaseService struct{}

func (s *releaseService) Get(ctx context.Context, id primitive.ObjectID) (*model.Release, error) {
	release := &model.Release{}
	if err := mongo.Repo.FindByID(ctx, release, id); err != nil {
		return nil, err
	}
	if release.DeletedAt != nil {
		return nil, mongoDriver.ErrNoDocuments
	}
	return release, nil
}

func (s *releaseService) updateStatus(ctx context.Context, releaseID primitive.ObjectID, status model.ReleaseStatus) error {
	filter := bson.M{
		"_id": releaseID,
		"status": bson.M{
			"$nin": []model.ReleaseStatus{model.ReleaseSucceeded, model.ReleaseFailed, model.ReleaseRolledBack, model.ReleaseSyncFailed, status},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	return mongo.Repo.UpdateOne(ctx, &model.Release{}, filter, update)
}

func (s *releaseService) UpdateStatus(ctx context.Context, releaseID primitive.ObjectID, status model.ReleaseStatus) error {
	return s.updateStatus(ctx, releaseID, status)
}

func (s *releaseService) UpdateStep(ctx context.Context, releaseID primitive.ObjectID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
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
		if err := s.createStepIfNotExists(ctx, releaseID, stepName, status, progress, message, start, end); err != nil {
			return err
		}

		nextSteps = append(nextSteps, model.ReleaseStep{
			Name:      stepName,
			Progress:  progress,
			Status:    status,
			Message:   message,
			StartTime: start,
			EndTime:   end,
		})
		return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
	}

	if currentStep.Status == model.StepFailed || currentStep.Status == model.StepSucceeded {
		return nil
	}

	update := bson.M{
		"steps.$.status":   status,
		"steps.$.progress": progress,
		"steps.$.message":  message,
		"updated_at":       time.Now(),
	}
	if start != nil {
		update["steps.$.start_time"] = *start
	}
	if end != nil {
		update["steps.$.end_time"] = *end
	}

	filter := bson.M{
		"_id": releaseID,
		"steps": bson.M{
			"$elemMatch": bson.M{
				"name": stepName,
				"status": bson.M{
					"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded},
				},
			},
		},
	}

	if err := mongo.Repo.UpdateOne(ctx, &model.Release{}, filter, bson.M{"$set": update}); err != nil {
		return err
	}

	applyReleaseStepUpdate(nextSteps, stepName, status, progress, message, start, end)
	return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
}

func (s *releaseService) createStepIfNotExists(ctx context.Context, releaseID primitive.ObjectID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	step := model.ReleaseStep{
		Name:      stepName,
		Progress:  progress,
		Status:    status,
		Message:   message,
		StartTime: start,
		EndTime:   end,
	}

	filter := bson.M{
		"_id": releaseID,
		"steps": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"name": stepName,
				},
			},
		},
	}

	update := bson.M{
		"$push": bson.M{
			"steps": step,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	return mongo.Repo.UpdateOne(ctx, &model.Release{}, filter, update)
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

func (s *releaseService) updateStatusFromSteps(ctx context.Context, releaseID primitive.ObjectID, releaseAction string, currentStatus model.ReleaseStatus, steps []model.ReleaseStep) error {
	nextStatus := model.DeriveReleaseStatusFromSteps(releaseAction, currentStatus, steps)
	if nextStatus == currentStatus {
		return nil
	}
	return s.updateStatus(ctx, releaseID, nextStatus)
}
