# API Spec

## Purpose

`devflow-verify-service` only exposes verify ingress APIs under `/api/v1/verify/*`.
These endpoints accept external execution facts and write them back through the approved verify path.

## Endpoint Groups

### Health
- `GET /api/v1/verify/healthz`

### Build status writeback
- `POST /api/v1/verify/tekton/events`
  - key fields: `manifest_id`, `status`, `pipeline_id`, `intent_id`, `external_ref`

### Build step writeback
- `POST /api/v1/verify/tekton/steps`
  - key fields: `manifest_id`, `pipeline_id`, `task_name`, `task_run`, `status`

### Release status writeback
- `POST /api/v1/verify/argo/events`
  - key fields: `release_id`, `status`, `intent_id`, `external_ref`

### Release step writeback
- `POST /api/v1/verify/release/steps`
  - key fields: `release_id`, `step_name`, `status`, `progress`

## Request Rules

- write endpoints use `X-Devflow-Verify-Token`
- if `VERIFY_SERVICE_SHARED_TOKEN` is set, write requests must pass token validation
- if `VERIFY_SERVICE_SHARED_TOKEN` is unset, local environments may access write endpoints without the token

## Response Rules

- health returns `200`
- write endpoints return normal success/error HTTP status codes from handler validation and writeback logic
- this repo does not expose pagination-based list APIs

## Error Rules

- request body missing / invalid ID / required field missing -> `400`
- shared token validation failed -> `401`
- Mongo update or internal writeback failure -> `500`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.
