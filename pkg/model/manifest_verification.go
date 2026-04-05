package model

import (
	"time"

	"github.com/google/uuid"
)

type ManifestVerification struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	ManifestID     uuid.UUID      `json:"manifest_id" db:"manifest_id"`
	IntentID       *uuid.UUID     `json:"intent_id,omitempty" db:"intent_id"`
	PipelineID     string         `json:"pipeline_id,omitempty" db:"pipeline_id"`
	Status         ManifestStatus `json:"status" db:"status"`
	ExternalRef    string         `json:"external_ref,omitempty" db:"external_ref"`
	Summary        string         `json:"summary,omitempty" db:"summary"`
	LastMessage    string         `json:"last_message,omitempty" db:"last_message"`
	Steps          []ManifestStep `json:"steps,omitempty" db:"steps"`
	Details        map[string]any `json:"details,omitempty" db:"details"`
	LastObservedAt time.Time      `json:"last_observed_at" db:"last_observed_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

func (ManifestVerification) CollectionName() string { return "manifest_verifications" }
