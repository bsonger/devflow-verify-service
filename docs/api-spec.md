# API Spec

## Purpose

`devflow-verify-service` defines the converged verify-ingress API surface under `/api/v1/verify/*`.
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
- verify-service does not expose public CRUD for `Manifest`, `Release`, or `Intent`
- write endpoints are for controller/observer execution facts only

## Response Rules

- health returns `200`
- write endpoints return normal success/error HTTP status codes from handler validation and writeback logic
- this repo does not expose pagination-based list APIs

## Error Rules

- request body missing / invalid ID / required field missing -> `400`
- shared token validation failed -> `401`
- internal writeback failure -> `500`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.

## Transitional Note

Generated Swagger artifacts in this repo still reflect the legacy Mongo/ObjectID handler layer and should be regenerated after the API layer adopts the relational contract.
