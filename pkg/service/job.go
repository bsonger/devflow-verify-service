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

var JobService = &jobService{}

type jobService struct{}

func (s *jobService) Get(ctx context.Context, id primitive.ObjectID) (*model.Job, error) {
	job := &model.Job{}
	if err := mongo.Repo.FindByID(ctx, job, id); err != nil {
		return nil, err
	}
	if job.DeletedAt != nil {
		return nil, mongoDriver.ErrNoDocuments
	}
	return job, nil
}

func (s *jobService) updateStatus(ctx context.Context, jobID primitive.ObjectID, status model.JobStatus) error {
	filter := bson.M{
		"_id": jobID,
		"status": bson.M{
			"$nin": []model.JobStatus{model.JobSucceeded, model.JobFailed, model.JobRolledBack, model.JobSyncFailed, status},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	return mongo.Repo.UpdateOne(ctx, &model.Job{}, filter, update)
}

func (s *jobService) UpdateStatus(ctx context.Context, jobID primitive.ObjectID, status model.JobStatus) error {
	return s.updateStatus(ctx, jobID, status)
}

func (s *jobService) UpdateStep(ctx context.Context, jobID primitive.ObjectID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	if stepName == "" {
		return nil
	}

	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	job, err := s.Get(ctx, jobID)
	if err != nil {
		return err
	}

	nextSteps := cloneJobSteps(job.Steps)
	currentStep := findJobStep(job.Steps, stepName)
	if currentStep == nil {
		if err := s.createStepIfNotExists(ctx, jobID, stepName, status, progress, message, start, end); err != nil {
			return err
		}

		nextSteps = append(nextSteps, model.JobStep{
			Name:      stepName,
			Progress:  progress,
			Status:    status,
			Message:   message,
			StartTime: start,
			EndTime:   end,
		})
		return s.updateStatusFromSteps(ctx, jobID, job.Type, job.Status, nextSteps)
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
		"_id": jobID,
		"steps": bson.M{
			"$elemMatch": bson.M{
				"name": stepName,
				"status": bson.M{
					"$nin": []model.StepStatus{model.StepFailed, model.StepSucceeded},
				},
			},
		},
	}

	if err := mongo.Repo.UpdateOne(ctx, &model.Job{}, filter, bson.M{"$set": update}); err != nil {
		return err
	}

	applyJobStepUpdate(nextSteps, stepName, status, progress, message, start, end)
	return s.updateStatusFromSteps(ctx, jobID, job.Type, job.Status, nextSteps)
}

func (s *jobService) createStepIfNotExists(ctx context.Context, jobID primitive.ObjectID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	step := model.JobStep{
		Name:      stepName,
		Progress:  progress,
		Status:    status,
		Message:   message,
		StartTime: start,
		EndTime:   end,
	}

	filter := bson.M{
		"_id": jobID,
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

	return mongo.Repo.UpdateOne(ctx, &model.Job{}, filter, update)
}

func findJobStep(steps []model.JobStep, stepName string) *model.JobStep {
	for _, step := range steps {
		if step.Name == stepName {
			current := step
			return &current
		}
	}
	return nil
}

func cloneJobSteps(steps []model.JobStep) []model.JobStep {
	if len(steps) == 0 {
		return nil
	}
	cloned := make([]model.JobStep, len(steps))
	copy(cloned, steps)
	return cloned
}

func applyJobStepUpdate(steps []model.JobStep, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) {
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

func (s *jobService) updateStatusFromSteps(ctx context.Context, jobID primitive.ObjectID, jobType string, currentStatus model.JobStatus, steps []model.JobStep) error {
	nextStatus := model.DeriveJobStatusFromSteps(jobType, currentStatus, steps)
	if nextStatus == currentStatus {
		return nil
	}
	return s.updateStatus(ctx, jobID, nextStatus)
}
