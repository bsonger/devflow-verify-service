package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

var ReleaseService = &releaseService{}

type releaseService struct{}

func (s *releaseService) Get(ctx context.Context, id uuid.UUID) (*releaseRecord, error) {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return nil, err
	}
	doc := &releaseDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, oid); err != nil {
		return nil, err
	}
	record := releaseRecordFromDoc(doc)
	if record.DeletedAt != nil {
		return nil, mongoDriver.ErrNoDocuments
	}
	return &record, nil
}

func (s *releaseService) updateStatus(ctx context.Context, releaseID uuid.UUID, status model.ReleaseStatus) error {
	oid, err := bridgeUUIDToObjectID(releaseID)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateOne(ctx, &releaseDoc{}, bson.M{
		"_id": oid,
		"status": bson.M{
			"$nin": []model.ReleaseStatus{model.ReleaseSucceeded, model.ReleaseFailed, model.ReleaseRolledBack, model.ReleaseSyncFailed, status},
		},
	}, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now()},
	})
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
		if err := s.createStepIfNotExists(ctx, releaseID, stepName, status, progress, message, start, end); err != nil {
			return err
		}
		nextSteps = append(nextSteps, model.ReleaseStep{Name: stepName, Progress: progress, Status: status, Message: message, StartTime: start, EndTime: end})
		return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
	}

	if currentStep.Status == model.StepFailed || currentStep.Status == model.StepSucceeded {
		return nil
	}

	oid, err := bridgeUUIDToObjectID(releaseID)
	if err != nil {
		return err
	}
	update := bson.M{"steps.$.status": status, "steps.$.progress": progress, "steps.$.message": message, "updated_at": time.Now()}
	if start != nil {
		update["steps.$.start_time"] = *start
	}
	if end != nil {
		update["steps.$.end_time"] = *end
	}

	if err := mongo.Repo.UpdateOne(ctx, &releaseDoc{}, bson.M{
		"_id": oid,
		"steps": bson.M{"$elemMatch": bson.M{
			"name":   stepName,
			"status": bson.M{"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded}},
		}},
	}, bson.M{"$set": update}); err != nil {
		return err
	}

	applyReleaseStepUpdate(nextSteps, stepName, status, progress, message, start, end)
	return s.updateStatusFromSteps(ctx, releaseID, release.Type, release.Status, nextSteps)
}

func (s *releaseService) createStepIfNotExists(ctx context.Context, releaseID uuid.UUID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	oid, err := bridgeUUIDToObjectID(releaseID)
	if err != nil {
		return err
	}
	step := model.ReleaseStep{Name: stepName, Progress: progress, Status: status, Message: message, StartTime: start, EndTime: end}
	return mongo.Repo.UpdateOne(ctx, &releaseDoc{}, bson.M{
		"_id":   oid,
		"steps": bson.M{"$not": bson.M{"$elemMatch": bson.M{"name": stepName}}},
	}, bson.M{
		"$push": bson.M{"steps": step},
		"$set":  bson.M{"updated_at": time.Now()},
	})
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
