package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestImageVerificationContract(t *testing.T) {
	typ := reflect.TypeOf(ImageVerification{})
	for _, field := range []string{"ImageID", "IntentID", "PipelineID", "Status", "Steps", "Details", "LastObservedAt"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("ImageVerification missing field %s", field)
		}
		if field == "ImageID" && f.Type != reflect.TypeOf(uuid.UUID{}) {
			t.Fatalf("ImageVerification.ImageID type = %v, want uuid.UUID", f.Type)
		}
	}
}

func TestReleaseVerificationContract(t *testing.T) {
	typ := reflect.TypeOf(ReleaseVerification{})
	for _, field := range []string{"ReleaseID", "IntentID", "Env", "Status", "Steps", "Details", "LastObservedAt"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("ReleaseVerification missing field %s", field)
		}
		if field == "ReleaseID" && f.Type != reflect.TypeOf(uuid.UUID{}) {
			t.Fatalf("ReleaseVerification.ReleaseID type = %v, want uuid.UUID", f.Type)
		}
	}
}

func TestBaseModelWithCreateDefault(t *testing.T) {
	var base BaseModel
	base.WithCreateDefault()

	if base.ID == uuid.Nil {
		t.Fatal("BaseModel.WithCreateDefault should assign a UUID")
	}
	if base.CreatedAt.IsZero() || base.UpdatedAt.IsZero() {
		t.Fatal("BaseModel.WithCreateDefault should set timestamps")
	}
}
