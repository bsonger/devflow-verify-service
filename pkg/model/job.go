package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus string

const (
	JobPending     JobStatus = "Pending"
	JobRunning     JobStatus = "Running"
	JobSucceeded   JobStatus = "Succeeded"
	JobFailed      JobStatus = "Failed"
	JobRollingBack JobStatus = "RollingBack"
	JobRolledBack  JobStatus = "RolledBack"
	JobSyncing     JobStatus = "Syncing"
	JobSyncFailed  JobStatus = "SyncFailed"

	JobInstall  string = "Install"
	JobUpgrade  string = "Upgrade"
	JobRollback string = "Rollback"
)

type Job struct {
	BaseModel `bson:",inline"`

	ExecutionIntentID *primitive.ObjectID `bson:"execution_intent_id,omitempty" json:"execution_intent_id,omitempty"`
	ManifestID        primitive.ObjectID  `bson:"manifest_id" json:"manifest_id"`
	Type              string              `bson:"type" json:"type"`
	Status            JobStatus           `bson:"status" json:"status"`
	Steps             []JobStep           `bson:"steps,omitempty" json:"steps,omitempty"`
}

func (Job) CollectionName() string { return "job" }

type JobStep struct {
	Name      string     `bson:"name" json:"name"`
	Progress  int32      `bson:"progress" json:"progress"`
	Status    StepStatus `bson:"status" json:"status"`
	Message   string     `bson:"message,omitempty" json:"message,omitempty"`
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
}

func DeriveJobStatusFromSteps(jobType string, currentStatus JobStatus, steps []JobStep) JobStatus {
	switch currentStatus {
	case JobSucceeded, JobFailed, JobRolledBack, JobSyncFailed:
		return currentStatus
	}

	if len(steps) == 0 {
		if currentStatus == "" {
			return JobPending
		}
		return currentStatus
	}

	allSucceeded := true
	anyFailed := false
	anyStarted := false

	for _, step := range steps {
		switch step.Status {
		case StepFailed:
			anyFailed = true
			allSucceeded = false
		case StepSucceeded:
			anyStarted = true
		case StepRunning:
			anyStarted = true
			allSucceeded = false
		default:
			allSucceeded = false
		}
	}

	if anyFailed {
		return JobFailed
	}
	if allSucceeded {
		if jobType == JobRollback {
			return JobRolledBack
		}
		return JobSucceeded
	}
	if anyStarted {
		return JobRunning
	}
	if currentStatus == "" {
		return JobPending
	}
	return currentStatus
}
