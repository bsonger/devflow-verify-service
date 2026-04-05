package service

import (
	"errors"
	"time"

	commonmodel "github.com/bsonger/devflow-common/model"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errInvalidUUIDBridge = errors.New("invalid bridged uuid")
	bridgeUUIDPrefix     = [4]byte{'d', 'f', 'l', 'w'}
)

type intentDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	Kind         string             `bson:"kind"`
	Status       string             `bson:"status"`
	ResourceType string             `bson:"resource_type"`
	ResourceID   primitive.ObjectID `bson:"resource_id"`
	ExternalRef  string             `bson:"external_ref,omitempty"`
	Message      string             `bson:"message,omitempty"`
	LastError    string             `bson:"last_error,omitempty"`
}

func (intentDoc) CollectionName() string { return "execution_intents" }

type manifestDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	PipelineID string               `bson:"pipeline_id,omitempty"`
	Steps      []model.ManifestStep `bson:"steps,omitempty"`
	Status     model.ManifestStatus `bson:"status"`
}

func (manifestDoc) CollectionName() string { return "manifests" }

type releaseDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	Type   string              `bson:"type"`
	Steps  []model.ReleaseStep `bson:"steps,omitempty"`
	Status model.ReleaseStatus `bson:"status"`
}

func (releaseDoc) CollectionName() string { return "releases" }

func bridgeObjectIDToUUID(id primitive.ObjectID) uuid.UUID {
	var raw [16]byte
	copy(raw[:4], bridgeUUIDPrefix[:])
	copy(raw[4:], id[:])
	return uuid.UUID(raw)
}

func bridgeUUIDToObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	raw := [16]byte(id)
	if raw[0] != bridgeUUIDPrefix[0] || raw[1] != bridgeUUIDPrefix[1] || raw[2] != bridgeUUIDPrefix[2] || raw[3] != bridgeUUIDPrefix[3] {
		return primitive.NilObjectID, errInvalidUUIDBridge
	}
	var oid primitive.ObjectID
	copy(oid[:], raw[4:])
	return oid, nil
}

func BridgeUUIDToObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	return bridgeUUIDToObjectID(id)
}

type releaseRecord struct {
	ID        uuid.UUID
	Status    model.ReleaseStatus
	Type      string
	Steps     []model.ReleaseStep
	DeletedAt *time.Time
}

type manifestRecord struct {
	ID         uuid.UUID
	PipelineID string
	Status     model.ManifestStatus
	Steps      []model.ManifestStep
	DeletedAt  *time.Time
}

func releaseRecordFromDoc(doc *releaseDoc) releaseRecord {
	return releaseRecord{
		ID:        bridgeObjectIDToUUID(doc.ID),
		Status:    doc.Status,
		Type:      doc.Type,
		Steps:     doc.Steps,
		DeletedAt: doc.DeletedAt,
	}
}

func manifestRecordFromDoc(doc *manifestDoc) manifestRecord {
	return manifestRecord{
		ID:         bridgeObjectIDToUUID(doc.ID),
		PipelineID: doc.PipelineID,
		Status:     doc.Status,
		Steps:      doc.Steps,
		DeletedAt:  doc.DeletedAt,
	}
}
