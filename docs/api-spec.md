# API Spec

## Purpose

`devflow-verify-service` defines the verify-ingress API surface under `/api/v1/verify/*`.
These endpoints accept external execution facts and write them back through the approved verify path.

## Swagger

- local UI: `/swagger/index.html`
- verify-scoped UI alias: `/api/v1/verify/swagger/index.html`
- generated source: `docs/generated/swagger/swagger.yaml`

## Endpoint Groups

### Health
- `GET /api/v1/verify/healthz`

### Build status writeback
- `POST /api/v1/verify/tekton/events`
  - key fields: `image_id`, `status`, `pipeline_id`, `intent_id`, `external_ref`

### Build step writeback
- `POST /api/v1/verify/tekton/steps`
  - key fields: `image_id`, `pipeline_id`, `task_name`, `task_run`, `status`

### Release status writeback
- `POST /api/v1/verify/argo/events`
  - key fields: `release_id`, `status`, `intent_id`, `external_ref`

### Release step writeback
- `POST /api/v1/verify/release/steps`
  - key fields: `release_id`, `step_name`, `status`, `progress`

## Request Rules

- write endpoints use `X-Devflow-Verify-Token`
- if `config.yaml` `auth.shared_token` is set, write requests must pass token validation
- if `auth.shared_token` is empty, local environments may access write endpoints without the token
- verify-service does not expose public CRUD for `Image`, `Release`, or `Intent`
- write endpoints are for controller and observer execution facts only

## Response Rules

- health returns `200` with `{ "data": { "service": "verify-service", "status": "ok" } }`
- write endpoints return `204 No Content` on success
- errors use `{ "error": { "code", "message", "details" } }`
- this repo does not expose pagination-based list APIs

## Error Rules

- request body missing, invalid ID, or required field missing -> `400 invalid_argument`
- missing target image or release -> `404 not_found`
- missing derived pipeline binding during step writeback -> `400 failed_precondition`
- shared token validation failed -> `401 unauthorized`
- internal writeback failure -> `500 internal`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.
