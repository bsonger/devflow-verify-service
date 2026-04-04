package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ManifestStatus string

const (
	ManifestPending   ManifestStatus = "Pending"
	ManifestRunning   ManifestStatus = "Running"
	ManifestSucceeded ManifestStatus = "Succeeded"
	ManifestFailed    ManifestStatus = "Failed"
)

type StepStatus string

const (
	StepPending   StepStatus = "Pending"
	StepRunning   StepStatus = "Running"
	StepSucceeded StepStatus = "Succeeded"
	StepFailed    StepStatus = "Failed"
)

type Manifest struct {
	BaseModel         `bson:",inline"`
	ExecutionIntentID *primitive.ObjectID `json:"execution_intent_id,omitempty" bson:"execution_intent_id,omitempty"`
	ApplicationId     primitive.ObjectID  `json:"application_id" bson:"application_id"`
	Name              string              `json:"name" bson:"name"`
	ApplicationName   string              `json:"application_name" bson:"application_name"`
	PipelineID        string              `json:"pipeline_id" bson:"pipeline_id"`
	Steps             []ManifestStep      `json:"steps" bson:"steps"`
	Status            ManifestStatus      `json:"status" bson:"status"`
}

type ManifestStep struct {
	TaskName  string     `bson:"task_name" json:"task_name"`
	TaskRun   string     `bson:"task_run,omitempty" json:"task_run,omitempty"`
	Status    StepStatus `bson:"status" json:"status"`
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Message   string     `bson:"message,omitempty" json:"message,omitempty"`
}

func (Manifest) CollectionName() string { return "manifests" }
