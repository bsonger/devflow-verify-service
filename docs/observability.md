# Observability

## Purpose

`devflow-verify-service` emits the shared backend telemetry baseline plus external execution-fact writeback context.

## Logs

Required structured fields:
- `resource`
- `resource_id`
- `image_id`
- `release_id`
- `intent_id`
- `pipeline_id`
- `task_name`
- `result`
- `error_code`

## Metrics

- use shared `devflow_http_*` ingress metrics
- avoid high-cardinality labels for pipeline/task identifiers
- rely on logs and spans, not label explosion, for detailed external callback debugging

## Tracing

- every business HTTP request should create a server span
- preserve inbound trace context when external controllers send it
- any downstream writeback or lookup must emit a client span with propagated trace context

## Health and readiness

- expose `GET /api/v1/verify/healthz`, `/metrics`, and readiness endpoints when the repo adds them
- exclude `/swagger/*` and diagnostics endpoints from business telemetry rollups

## Failure modes

Watch for:
- shared token validation failures
- missing target resource or pipeline binding failures
- external callback payload drift
- writeback convergence failures into release-service state

## Dashboards and runbooks

Use the shared backend dashboard/runbook set plus verify-specific callback views when they exist.
