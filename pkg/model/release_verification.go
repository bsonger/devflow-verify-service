package model

import (
	"time"

	"github.com/google/uuid"
)

type ReleaseVerification struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	ReleaseID      uuid.UUID      `json:"release_id" db:"release_id"`
	IntentID       *uuid.UUID     `json:"intent_id,omitempty" db:"intent_id"`
	Env            string         `json:"env,omitempty" db:"env"`
	Status         ReleaseStatus  `json:"status" db:"status"`
	ExternalRef    string         `json:"external_ref,omitempty" db:"external_ref"`
	Summary        string         `json:"summary,omitempty" db:"summary"`
	LastMessage    string         `json:"last_message,omitempty" db:"last_message"`
	Steps          []ReleaseStep  `json:"steps,omitempty" db:"steps"`
	Details        map[string]any `json:"details,omitempty" db:"details"`
	LastObservedAt time.Time      `json:"last_observed_at" db:"last_observed_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

func (ReleaseVerification) CollectionName() string { return "release_verifications" }
