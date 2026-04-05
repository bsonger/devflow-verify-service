package model

import "time"

type StepStatus string
type ManifestStatus string
type ReleaseStatus string

const (
	StepPending   StepStatus = "Pending"
	StepRunning   StepStatus = "Running"
	StepSucceeded StepStatus = "Succeeded"
	StepFailed    StepStatus = "Failed"
)

const (
	ManifestPending   ManifestStatus = "Pending"
	ManifestRunning   ManifestStatus = "Running"
	ManifestSucceeded ManifestStatus = "Succeeded"
	ManifestFailed    ManifestStatus = "Failed"
)

const (
	ReleasePending     ReleaseStatus = "Pending"
	ReleaseRunning     ReleaseStatus = "Running"
	ReleaseSucceeded   ReleaseStatus = "Succeeded"
	ReleaseFailed      ReleaseStatus = "Failed"
	ReleaseRollingBack ReleaseStatus = "RollingBack"
	ReleaseRolledBack  ReleaseStatus = "RolledBack"
	ReleaseSyncing     ReleaseStatus = "Syncing"
	ReleaseSyncFailed  ReleaseStatus = "SyncFailed"
)

type ManifestStep struct {
	TaskName  string     `json:"task_name"`
	TaskRun   string     `json:"task_run,omitempty"`
	Status    StepStatus `json:"status"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type ReleaseStep struct {
	Name      string     `json:"name"`
	Progress  int32      `json:"progress"`
	Status    StepStatus `json:"status"`
	Message   string     `json:"message,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}
