package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReleaseStatus string

const (
	ReleasePending     ReleaseStatus = "Pending"
	ReleaseRunning     ReleaseStatus = "Running"
	ReleaseSucceeded   ReleaseStatus = "Succeeded"
	ReleaseFailed      ReleaseStatus = "Failed"
	ReleaseRollingBack ReleaseStatus = "RollingBack"
	ReleaseRolledBack  ReleaseStatus = "RolledBack"
	ReleaseSyncing     ReleaseStatus = "Syncing"
	ReleaseSyncFailed  ReleaseStatus = "SyncFailed"

	ReleaseInstall  string = "Install"
	ReleaseUpgrade  string = "Upgrade"
	ReleaseRollback string = "Rollback"
)

type Release struct {
	BaseModel `bson:",inline"`

	ExecutionIntentID *primitive.ObjectID `bson:"execution_intent_id,omitempty" json:"execution_intent_id,omitempty"`
	ManifestID        primitive.ObjectID  `bson:"manifest_id" json:"manifest_id"`
	Type              string              `bson:"type" json:"type"`
	Status            ReleaseStatus       `bson:"status" json:"status"`
	Steps             []ReleaseStep       `bson:"steps,omitempty" json:"steps,omitempty"`
}

func (Release) CollectionName() string { return "release" }

type ReleaseStep struct {
	Name      string     `bson:"name" json:"name"`
	Progress  int32      `bson:"progress" json:"progress"`
	Status    StepStatus `bson:"status" json:"status"`
	Message   string     `bson:"message,omitempty" json:"message,omitempty"`
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
}

func DeriveReleaseStatusFromSteps(releaseAction string, currentStatus ReleaseStatus, steps []ReleaseStep) ReleaseStatus {
	switch currentStatus {
	case ReleaseSucceeded, ReleaseFailed, ReleaseRolledBack, ReleaseSyncFailed:
		return currentStatus
	}

	if len(steps) == 0 {
		if currentStatus == "" {
			return ReleasePending
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
		return ReleaseFailed
	}
	if allSucceeded {
		if releaseAction == ReleaseRollback {
			return ReleaseRolledBack
		}
		return ReleaseSucceeded
	}
	if anyStarted {
		return ReleaseRunning
	}
	if currentStatus == "" {
		return ReleasePending
	}
	return currentStatus
}
