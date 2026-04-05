package api

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bsonger/devflow-service-common/httpx"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/bsonger/devflow-verify-service/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var VerifyRouteApi = NewVerifyHandler()

type VerifyHandler struct {
	manifestSvc manifestWriteService
	releaseSvc  releaseWriteService
	intentSvc   intentWriteService
}

type manifestWriteService interface {
	AssignPipelineID(ctx context.Context, manifestID uuid.UUID, pipelineID string) error
	UpdateManifestStatusByID(ctx context.Context, manifestID uuid.UUID, status model.ManifestStatus) error
	UpdateStepStatus(ctx context.Context, pipelineID, taskName string, status model.StepStatus, message string, start, end *time.Time) error
	BindTaskRun(ctx context.Context, pipelineID, taskName, taskRun string) error
}

type releaseWriteService interface {
	UpdateStatus(ctx context.Context, releaseID uuid.UUID, status model.ReleaseStatus) error
	UpdateStep(ctx context.Context, releaseID uuid.UUID, stepName string, status model.StepStatus, progress int32, message string, start, end *time.Time) error
}

type intentWriteService interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, externalRef, message string) error
	UpdateStatusByResource(ctx context.Context, kind string, resourceID uuid.UUID, status string, externalRef, message string) error
}

var loadManifestPipelineID = func(ctx context.Context, id uuid.UUID) (string, error) {
	manifest, err := service.ManifestService.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return manifest.PipelineID, nil
}

const VerifyTokenHeader = "X-Devflow-Verify-Token"

type VerifyBuildStatusRequest struct {
	IntentID    string               `json:"intent_id,omitempty"`
	ManifestID  string               `json:"manifest_id" binding:"required"`
	PipelineID  string               `json:"pipeline_id,omitempty"`
	Status      model.ManifestStatus `json:"status" binding:"required"`
	Message     string               `json:"message,omitempty"`
	ExternalRef string               `json:"external_ref,omitempty"`
}

type VerifyReleaseStatusRequest struct {
	IntentID    string              `json:"intent_id,omitempty"`
	ReleaseID   string              `json:"release_id" binding:"required"`
	Status      model.ReleaseStatus `json:"status" binding:"required"`
	Message     string              `json:"message,omitempty"`
	ExternalRef string              `json:"external_ref,omitempty"`
}

type VerifyReleaseStepRequest struct {
	ReleaseID string           `json:"release_id" binding:"required"`
	StepName  string           `json:"step_name" binding:"required"`
	Status    model.StepStatus `json:"status" binding:"required"`
	Progress  int32            `json:"progress,omitempty"`
	Message   string           `json:"message,omitempty"`
	StartTime *time.Time       `json:"start_time,omitempty"`
	EndTime   *time.Time       `json:"end_time,omitempty"`
}

type VerifyBuildStepRequest struct {
	ManifestID string           `json:"manifest_id" binding:"required"`
	PipelineID string           `json:"pipeline_id,omitempty"`
	TaskName   string           `json:"task_name" binding:"required"`
	TaskRun    string           `json:"task_run,omitempty"`
	Status     model.StepStatus `json:"status" binding:"required"`
	Message    string           `json:"message,omitempty"`
	StartTime  *time.Time       `json:"start_time,omitempty"`
	EndTime    *time.Time       `json:"end_time,omitempty"`
}

func NewVerifyHandler() *VerifyHandler {
	return &VerifyHandler{
		manifestSvc: service.ManifestService,
		releaseSvc:  service.ReleaseService,
		intentSvc:   service.IntentService,
	}
}

func RequireVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := strings.TrimSpace(os.Getenv("VERIFY_SERVICE_SHARED_TOKEN"))
		if expected == "" {
			c.Next()
			return
		}

		token := strings.TrimSpace(c.GetHeader(VerifyTokenHeader))
		if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
			httpx.WriteError(c, http.StatusUnauthorized, "unauthorized", "unauthorized", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Health
// @Summary Verify Service 健康检查
// @Tags Verify
// @Success 200 {object} httpx.DataResponse[map[string]string]
// @Router /api/v1/verify/healthz [get]
func (h *VerifyHandler) Health(c *gin.Context) {
	httpx.WriteData(c, http.StatusOK, gin.H{
		"service": "verify-service",
		"status":  "ok",
	})
}

// HandleArgoEvent
// @Summary 回写发布状态
// @Description 由 Argo 或外部发布观察器回写 Release 级状态
// @Tags Verify
// @Accept json
// @Produce json
// @Param data body VerifyReleaseStatusRequest true "Release Status Data"
// @Success 204
// @Router /api/v1/verify/argo/events [post]
func (h *VerifyHandler) HandleArgoEvent(c *gin.Context) {
	var req VerifyReleaseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	releaseID, err := uuid.Parse(req.ReleaseID)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid release_id", nil)
		return
	}

	if err := h.releaseSvc.UpdateStatus(c.Request.Context(), releaseID, req.Status); err != nil {
		writeVerifyError(c, err)
		return
	}

	if req.IntentID != "" {
		if intentID, err := uuid.Parse(req.IntentID); err == nil {
			_ = h.intentSvc.UpdateStatus(c.Request.Context(), intentID, mapReleaseStatusToIntentStatus(req.Status), req.ExternalRef, req.Message)
		}
	} else {
		_ = h.intentSvc.UpdateStatusByResource(c.Request.Context(), "release", releaseID, mapReleaseStatusToIntentStatus(req.Status), req.ExternalRef, req.Message)
	}

	httpx.WriteNoContent(c)
}

// HandleTektonEvent
// @Summary 回写构建状态
// @Description 由 Tekton 或外部构建观察器回写 Manifest 级状态
// @Tags Verify
// @Accept json
// @Produce json
// @Param data body VerifyBuildStatusRequest true "Build Status Data"
// @Success 204
// @Router /api/v1/verify/tekton/events [post]
func (h *VerifyHandler) HandleTektonEvent(c *gin.Context) {
	var req VerifyBuildStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	manifestID, err := uuid.Parse(req.ManifestID)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid manifest_id", nil)
		return
	}

	if req.PipelineID != "" {
		if err := h.manifestSvc.AssignPipelineID(c.Request.Context(), manifestID, req.PipelineID); err != nil {
			writeVerifyError(c, err)
			return
		}
	}

	if err := h.manifestSvc.UpdateManifestStatusByID(c.Request.Context(), manifestID, req.Status); err != nil {
		writeVerifyError(c, err)
		return
	}

	if req.IntentID != "" {
		if intentID, err := uuid.Parse(req.IntentID); err == nil {
			_ = h.intentSvc.UpdateStatus(c.Request.Context(), intentID, mapManifestStatusToIntentStatus(req.Status), req.ExternalRef, req.Message)
		}
	} else {
		_ = h.intentSvc.UpdateStatusByResource(c.Request.Context(), "build", manifestID, mapManifestStatusToIntentStatus(req.Status), req.ExternalRef, req.Message)
	}

	httpx.WriteNoContent(c)
}

// HandleTektonStepEvent
// @Summary 回写构建步骤
// @Description 由 Tekton TaskRun 观察器回写 Manifest steps
// @Tags Verify
// @Accept json
// @Produce json
// @Param data body VerifyBuildStepRequest true "Build Step Data"
// @Success 204
// @Router /api/v1/verify/tekton/steps [post]
func (h *VerifyHandler) HandleTektonStepEvent(c *gin.Context) {
	var req VerifyBuildStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	manifestID, err := uuid.Parse(req.ManifestID)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid manifest_id", nil)
		return
	}

	if req.PipelineID == "" {
		pipelineID, err := loadManifestPipelineID(c.Request.Context(), manifestID)
		if err != nil {
			writeVerifyError(c, err)
			return
		}
		req.PipelineID = pipelineID
	}

	if req.PipelineID == "" {
		httpx.WriteError(c, http.StatusBadRequest, "failed_precondition", "pipeline_id is required until manifest is bound", nil)
		return
	}

	if req.TaskRun != "" {
		if err := h.manifestSvc.BindTaskRun(c.Request.Context(), req.PipelineID, req.TaskName, req.TaskRun); err != nil {
			writeVerifyError(c, err)
			return
		}
	}

	if err := h.manifestSvc.UpdateStepStatus(c.Request.Context(), req.PipelineID, req.TaskName, req.Status, req.Message, req.StartTime, req.EndTime); err != nil {
		writeVerifyError(c, err)
		return
	}

	httpx.WriteNoContent(c)
}

// HandleReleaseStepEvent
// @Summary 回写发布步骤
// @Description 由 Argo Application / Deployment / Rollout 观察器回写 Release steps
// @Tags Verify
// @Accept json
// @Produce json
// @Param data body VerifyReleaseStepRequest true "Release Step Data"
// @Success 204
// @Router /api/v1/verify/release/steps [post]
func (h *VerifyHandler) HandleReleaseStepEvent(c *gin.Context) {
	var req VerifyReleaseStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	releaseID, err := uuid.Parse(req.ReleaseID)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid release_id", nil)
		return
	}

	if err := h.releaseSvc.UpdateStep(c.Request.Context(), releaseID, req.StepName, req.Status, req.Progress, req.Message, req.StartTime, req.EndTime); err != nil {
		writeVerifyError(c, err)
		return
	}

	httpx.WriteNoContent(c)
}

func writeVerifyError(c *gin.Context, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
		return
	}
	httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
}

func mapManifestStatusToIntentStatus(status model.ManifestStatus) string {
	switch status {
	case model.ManifestSucceeded:
		return "Succeeded"
	case model.ManifestFailed:
		return "Failed"
	default:
		return "Running"
	}
}

func mapReleaseStatusToIntentStatus(status model.ReleaseStatus) string {
	switch status {
	case model.ReleaseSucceeded, model.ReleaseRolledBack:
		return "Succeeded"
	case model.ReleaseFailed, model.ReleaseSyncFailed:
		return "Failed"
	default:
		return "Running"
	}
}
