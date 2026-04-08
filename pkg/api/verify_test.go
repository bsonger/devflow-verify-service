package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubImageWriteService struct {
	assignPipelineIDFn    func(context.Context, uuid.UUID, string) error
	updateImageStatusByID func(context.Context, uuid.UUID, model.ImageStatus) error
	updateStepStatusFn    func(context.Context, string, string, model.StepStatus, string, *time.Time, *time.Time) error
	bindTaskRunFn         func(context.Context, string, string, string) error
}

func (s stubImageWriteService) AssignPipelineID(ctx context.Context, imageID uuid.UUID, pipelineID string) error {
	return s.assignPipelineIDFn(ctx, imageID, pipelineID)
}

func (s stubImageWriteService) UpdateImageStatusByID(ctx context.Context, imageID uuid.UUID, status model.ImageStatus) error {
	return s.updateImageStatusByID(ctx, imageID, status)
}

func (s stubImageWriteService) UpdateStepStatus(ctx context.Context, pipelineID, taskName string, status model.StepStatus, message string, start, end *time.Time) error {
	return s.updateStepStatusFn(ctx, pipelineID, taskName, status, message, start, end)
}

func (s stubImageWriteService) BindTaskRun(ctx context.Context, pipelineID, taskName, taskRun string) error {
	return s.bindTaskRunFn(ctx, pipelineID, taskName, taskRun)
}

type stubReleaseWriteService struct {
	updateStatusFn func(context.Context, uuid.UUID, model.ReleaseStatus) error
	updateStepFn   func(context.Context, uuid.UUID, string, model.StepStatus, int32, string, *time.Time, *time.Time) error
}

func (s stubReleaseWriteService) UpdateStatus(ctx context.Context, releaseID uuid.UUID, status model.ReleaseStatus) error {
	return s.updateStatusFn(ctx, releaseID, status)
}

func (s stubReleaseWriteService) UpdateStep(ctx context.Context, releaseID uuid.UUID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error {
	return s.updateStepFn(ctx, releaseID, stepName, status, progress, message, start, end)
}

type stubIntentWriteService struct {
	updateStatusFn           func(context.Context, uuid.UUID, string, string, string) error
	updateStatusByResourceFn func(context.Context, string, uuid.UUID, string, string, string) error
}

func (s stubIntentWriteService) UpdateStatus(ctx context.Context, id uuid.UUID, status string, externalRef, message string) error {
	return s.updateStatusFn(ctx, id, status, externalRef, message)
}

func (s stubIntentWriteService) UpdateStatusByResource(ctx context.Context, kind string, resourceID uuid.UUID, status string, externalRef, message string) error {
	return s.updateStatusByResourceFn(ctx, kind, resourceID, status, externalRef, message)
}

func TestRequireVerifyTokenAllowsWhenUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("VERIFY_SERVICE_SHARED_TOKEN", "")

	router := gin.New()
	router.Use(RequireVerifyToken())
	router.POST("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestRequireVerifyTokenRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("VERIFY_SERVICE_SHARED_TOKEN", "secret-token")

	router := gin.New()
	router.Use(RequireVerifyToken())
	router.POST("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set(VerifyTokenHeader, "wrong-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusUnauthorized)
	}

	var payload struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Error.Code != "unauthorized" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
}

func TestRequireVerifyTokenAcceptsValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("VERIFY_SERVICE_SHARED_TOKEN", "secret-token")

	router := gin.New()
	router.Use(RequireVerifyToken())
	router.POST("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set(VerifyTokenHeader, "secret-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestRequireVerifyTokenUsesEnvironmentAtRequestTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequireVerifyToken())
	router.POST("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	if err := os.Setenv("VERIFY_SERVICE_SHARED_TOKEN", "dynamic-secret"); err != nil {
		t.Fatalf("setenv failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv("VERIFY_SERVICE_SHARED_TOKEN")
	})

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set(VerifyTokenHeader, "dynamic-secret")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestHealthReturnsDataEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/verify/healthz", VerifyRouteApi.Health)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/verify/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Data["status"] != "ok" {
		t.Fatalf("unexpected payload: %#v", payload.Data)
	}
}

func TestHandleArgoEventReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &VerifyHandler{
		imageSvc:    stubImageWriteService{},
		releaseSvc: stubReleaseWriteService{
			updateStatusFn: func(context.Context, uuid.UUID, model.ReleaseStatus) error { return nil },
			updateStepFn: func(context.Context, uuid.UUID, string, model.StepStatus, int32, string, *time.Time, *time.Time) error {
				return nil
			},
		},
		intentSvc: stubIntentWriteService{
			updateStatusFn:           func(context.Context, uuid.UUID, string, string, string) error { return nil },
			updateStatusByResourceFn: func(context.Context, string, uuid.UUID, string, string, string) error { return nil },
		},
	}

	r := gin.New()
	r.POST("/api/v1/verify/argo/events", handler.HandleArgoEvent)

	body := bytes.NewBufferString(`{"release_id":"11111111-1111-1111-1111-111111111111","status":"Succeeded"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/verify/argo/events", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestHandleTektonStepEventMissingPipelineReturnsFailedPrecondition(t *testing.T) {
	gin.SetMode(gin.TestMode)
	origLoad := loadImagePipelineID
	loadImagePipelineID = func(context.Context, uuid.UUID) (string, error) { return "", nil }
	t.Cleanup(func() { loadImagePipelineID = origLoad })

	handler := &VerifyHandler{
		imageSvc: stubImageWriteService{
			assignPipelineIDFn:    func(context.Context, uuid.UUID, string) error { return nil },
			updateImageStatusByID: func(context.Context, uuid.UUID, model.ImageStatus) error { return nil },
			updateStepStatusFn: func(context.Context, string, string, model.StepStatus, string, *time.Time, *time.Time) error {
				return nil
			},
			bindTaskRunFn: func(context.Context, string, string, string) error { return nil },
		},
		releaseSvc: stubReleaseWriteService{},
		intentSvc:  stubIntentWriteService{},
	}

	r := gin.New()
	r.POST("/api/v1/verify/tekton/steps", handler.HandleTektonStepEvent)

	body := bytes.NewBufferString(`{"image_id":"11111111-1111-1111-1111-111111111111","task_name":"build","status":"Running"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/verify/tekton/steps", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d want %d", rec.Code, http.StatusBadRequest)
	}
	var payload struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Error.Code != "failed_precondition" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
}

func TestHandleReleaseStepEventNotFoundReturnsErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &VerifyHandler{
		imageSvc:    stubImageWriteService{},
		releaseSvc: stubReleaseWriteService{
			updateStatusFn: func(context.Context, uuid.UUID, model.ReleaseStatus) error { return nil },
			updateStepFn: func(context.Context, uuid.UUID, string, model.StepStatus, int32, string, *time.Time, *time.Time) error {
				return sql.ErrNoRows
			},
		},
		intentSvc: stubIntentWriteService{},
	}

	r := gin.New()
	r.POST("/api/v1/verify/release/steps", handler.HandleReleaseStepEvent)

	body := bytes.NewBufferString(`{"release_id":"11111111-1111-1111-1111-111111111111","step_name":"deploy","status":"Running"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/verify/release/steps", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d want %d", rec.Code, http.StatusNotFound)
	}
}
