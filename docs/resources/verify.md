# Verify

## Ownership

- owner repo: `devflow-verify-service`
- authoritative handler/API file: `pkg/api/verify.go`
- authoritative API doc: `docs/api-spec.md`
- swagger source: `docs/generated/swagger/swagger.yaml`

## Purpose

`Verify` 不是传统的 CRUD 资源，而是一组写回入口，用于接收外部执行事实并更新 `Manifest` / `Release` / `Intent`。

## Public endpoints

- `GET /api/v1/verify/healthz`
- `POST /api/v1/verify/argo/events`
- `POST /api/v1/verify/tekton/events`
- `POST /api/v1/verify/tekton/steps`
- `POST /api/v1/verify/release/steps`

## Auth

写接口使用请求头：
- `X-Devflow-Verify-Token`

若设置了 `VERIFY_SERVICE_SHARED_TOKEN`，则写接口必须通过鉴权。

## Request payloads

### `VerifyBuildStatusRequest`
Used by:
- `POST /api/v1/verify/tekton/events`

| Field | Type | Required | Description |
|---|---|---|---|
| `intent_id` | `string` | optional | 对应 build intent |
| `image_id` | `string` | required | 目标 manifest ID |
| `pipeline_id` | `string` | optional | pipeline 标识；若提供会回写到 manifest |
| `status` | `ManifestStatus` | required | Manifest 状态 |
| `message` | `string` | optional | 状态说明 |
| `external_ref` | `string` | optional | 外部引用 |

### `VerifyBuildStepRequest`
Used by:
- `POST /api/v1/verify/tekton/steps`

| Field | Type | Required | Description |
|---|---|---|---|
| `image_id` | `string` | required | 目标 manifest ID |
| `pipeline_id` | `string` | optional | pipeline 标识；为空时服务会尝试从 manifest 读取 |
| `task_name` | `string` | required | Tekton task 名 |
| `task_run` | `string` | optional | TaskRun 名 |
| `status` | `StepStatus` | required | step 状态 |
| `message` | `string` | optional | step 消息 |
| `start_time` | `*time.Time` | optional | 开始时间 |
| `end_time` | `*time.Time` | optional | 结束时间 |

Special rule:
- 如果 `pipeline_id` 为空且 manifest 上也未绑定 pipeline，接口返回错误：`pipeline_id is required until manifest is bound`

### `VerifyReleaseStatusRequest`
Used by:
- `POST /api/v1/verify/argo/events`

| Field | Type | Required | Description |
|---|---|---|---|
| `intent_id` | `string` | optional | 对应 release intent |
| `release_id` | `string` | required | 目标 release ID |
| `status` | `ReleaseStatus` | required | Release 状态 |
| `message` | `string` | optional | 状态说明 |
| `external_ref` | `string` | optional | 外部引用 |

### `VerifyReleaseStepRequest`
Used by:
- `POST /api/v1/verify/release/steps`

| Field | Type | Required | Description |
|---|---|---|---|
| `release_id` | `string` | required | 目标 release ID |
| `step_name` | `string` | required | 发布步骤名 |
| `status` | `StepStatus` | required | 步骤状态 |
| `progress` | `int32` | optional | 进度 |
| `message` | `string` | optional | 步骤消息 |
| `start_time` | `*time.Time` | optional | 开始时间 |
| `end_time` | `*time.Time` | optional | 结束时间 |

## Enums used by Verify

### `ManifestStatus`
- `Pending`
- `Running`
- `Succeeded`
- `Failed`

### `ReleaseStatus`
- `Pending`
- `Running`
- `Succeeded`
- `Failed`
- `RollingBack`
- `RolledBack`
- `Syncing`
- `SyncFailed`

### `StepStatus`
- `Pending`
- `Running`
- `Succeeded`
- `Failed`

## Writeback targets

- build status / step -> `Manifest`
- release status / step -> `Release`
- status convergence may also update related `Intent`

## Validation notes

- `image_id` / `release_id` 必须是合法 UUID
- required 字段由 handler 的 `binding:"required"` 定义
- pipeline/task/release step 语义还受 service 写回逻辑约束

## Source pointers

- router: `pkg/router/verify.go`
- handler: `pkg/api/verify.go`
- manifest writeback: `pkg/service/image.go`
- release writeback: `pkg/service/release.go`
- intent writeback: `pkg/service/intent.go`
